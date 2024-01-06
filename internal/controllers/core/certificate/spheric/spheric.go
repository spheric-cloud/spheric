// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package spheric

import (
	"spheric.cloud/spheric/internal/controllers/core/certificate/compute"
	"spheric.cloud/spheric/internal/controllers/core/certificate/generic"
	"spheric.cloud/spheric/internal/controllers/core/certificate/networking"
	"spheric.cloud/spheric/internal/controllers/core/certificate/storage"
)

var Recognizers []generic.CertificateSigningRequestRecognizer

func init() {
	Recognizers = append(Recognizers, compute.Recognizers...)
	Recognizers = append(Recognizers, storage.Recognizers...)
	Recognizers = append(Recognizers, networking.Recognizers...)
}
