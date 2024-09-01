// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core_test

import (
	"crypto/x509"
	"crypto/x509/pkix"

	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	certificatesv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"
	utilcertificate "spheric.cloud/spheric/utils/certificate"
)

var _ = Describe("CertificateApprovalController", func() {
	DescribeTable("certificate approval",
		func(ctx SpecContext, commonName, organization string) {
			By("creating a certificate signing request")
			csr, _, _, err := utilcertificate.GenerateAndCreateCertificateSigningRequest(
				ctx,
				k8sClient,
				certificatesv1.KubeAPIServerClientSignerName,
				&x509.CertificateRequest{
					Subject: pkix.Name{
						CommonName:   commonName,
						Organization: []string{organization},
					},
				},
				utilcertificate.DefaultKubeAPIServerClientGetUsages,
				nil,
			)
			Expect(err).NotTo(HaveOccurred())

			By("waiting for the csr to be approved and a certificate to be available")
			Eventually(ctx, Object(csr)).Should(
				HaveField("Status.Conditions", ContainElement(SatisfyAll(
					HaveField("Type", certificatesv1.CertificateApproved),
					HaveField("Status", corev1.ConditionTrue),
				))),
			)
		},
		Entry("fleet",
			corev1alpha1.FleetCommonName("my-pool"),
			corev1alpha1.FleetsGroup,
		),
	)
})
