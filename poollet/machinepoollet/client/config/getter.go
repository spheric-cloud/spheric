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
	computev1alpha1 "spheric.cloud/spheric/api/compute/v1alpha1"
	utilcertificate "spheric.cloud/spheric/utils/certificate"
	"spheric.cloud/spheric/utils/client/config"
)

var log = ctrl.Log.WithName("client").WithName("config")

func NewGetter(machinePoolName string) (*config.Getter, error) {
	return config.NewGetter(config.GetterOptions{
		Name:       "machinepoollet",
		SignerName: certificatesv1.KubeAPIServerClientSignerName,
		Template: &x509.CertificateRequest{
			Subject: pkix.Name{
				CommonName:   computev1alpha1.MachinePoolCommonName(machinePoolName),
				Organization: []string{computev1alpha1.MachinePoolsGroup},
			},
		},
		GetUsages:      utilcertificate.DefaultKubeAPIServerClientGetUsages,
		NetworkContext: egressselector.ControlPlane.AsNetworkContext(),
	})
}

func NewGetterOrDie(machinePoolName string) *config.Getter {
	getter, err := NewGetter(machinePoolName)
	if err != nil {
		log.Error(err, "Error creating getter")
		os.Exit(1)
	}
	return getter
}
