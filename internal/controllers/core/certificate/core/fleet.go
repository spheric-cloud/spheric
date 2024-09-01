// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"crypto/x509"
	"fmt"
	"strings"

	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"

	"golang.org/x/exp/slices"
	authv1 "k8s.io/api/authorization/v1"
	certificatesv1 "k8s.io/api/certificates/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"spheric.cloud/spheric/internal/controllers/core/certificate/generic"
)

var (
	FleetRequiredUsages = sets.New[certificatesv1.KeyUsage](
		certificatesv1.UsageDigitalSignature,
		certificatesv1.UsageKeyEncipherment,
		certificatesv1.UsageClientAuth,
	)
)

func IsFleetClientCert(csr *certificatesv1.CertificateSigningRequest, x509cr *x509.CertificateRequest) bool {
	if csr.Spec.SignerName != certificatesv1.KubeAPIServerClientSignerName {
		return false
	}

	return ValidateFleetClientCSR(x509cr, sets.New(csr.Spec.Usages...)) == nil
}

func ValidateFleetClientCSR(req *x509.CertificateRequest, usages sets.Set[certificatesv1.KeyUsage]) error {
	if !slices.Equal([]string{corev1alpha1.FleetsGroup}, req.Subject.Organization) {
		return fmt.Errorf("organization is not %s", corev1alpha1.FleetsGroup)
	}

	if len(req.DNSNames) > 0 {
		return fmt.Errorf("dns subject alternative names are not allowed")
	}
	if len(req.EmailAddresses) > 0 {
		return fmt.Errorf("email subject alternative names are not allowed")
	}
	if len(req.IPAddresses) > 0 {
		return fmt.Errorf("ip subject alternative names are not allowed")
	}
	if len(req.URIs) > 0 {
		return fmt.Errorf("uri subject alternative names are not allowed")
	}

	if !strings.HasPrefix(req.Subject.CommonName, corev1alpha1.FleetUserNamePrefix) {
		return fmt.Errorf("subject common name does not begin with %s", corev1alpha1.FleetUserNamePrefix)
	}

	if !FleetRequiredUsages.Equal(usages) {
		return fmt.Errorf("usages did not match %v", sets.List(FleetRequiredUsages))
	}

	return nil
}

var (
	FleetRecognizer = generic.NewCertificateSigningRequestRecognizer(
		IsFleetClientCert,
		authv1.ResourceAttributes{
			Group:       certificatesv1.GroupName,
			Resource:    "certificatesigningrequests",
			Verb:        "create",
			Subresource: "fleetclient",
		},
		"Auto approving fleet client certificate after SubjectAccessReview.",
	)
)

func init() {
	Recognizers = append(Recognizers, FleetRecognizer)
}
