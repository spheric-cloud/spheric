// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"os"

	certificatesv1 "k8s.io/api/certificates/v1"
	"k8s.io/apiserver/pkg/server/egressselector"
	ctrl "sigs.k8s.io/controller-runtime"
	storagev1alpha1 "spheric.cloud/spheric/api/storage/v1alpha1"
	utilcertificate "spheric.cloud/spheric/utils/certificate"
	"spheric.cloud/spheric/utils/client/config"
)

var log = ctrl.Log.WithName("client").WithName("config")

func NewGetter(bucketPoolName string) (*config.Getter, error) {
	return config.NewGetter(config.GetterOptions{
		Name:       "bucketpoollet",
		SignerName: certificatesv1.KubeAPIServerClientSignerName,
		Template: &x509.CertificateRequest{
			Subject: pkix.Name{
				CommonName:   storagev1alpha1.BucketPoolCommonName(bucketPoolName),
				Organization: []string{storagev1alpha1.BucketPoolsGroup},
			},
		},
		GetUsages:      utilcertificate.DefaultKubeAPIServerClientGetUsages,
		NetworkContext: egressselector.ControlPlane.AsNetworkContext(),
	})
}

func NewGetterOrDie(bucketPoolName string) *config.Getter {
	getter, err := NewGetter(bucketPoolName)
	if err != nil {
		log.Error(err, "Error creating getter")
		os.Exit(1)
	}
	return getter
}
