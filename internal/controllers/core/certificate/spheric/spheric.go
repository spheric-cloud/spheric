// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package spheric

import (
	"spheric.cloud/spheric/internal/controllers/core/certificate/core"
	"spheric.cloud/spheric/internal/controllers/core/certificate/generic"
)

var Recognizers []generic.CertificateSigningRequestRecognizer

func init() {
	Recognizers = append(Recognizers, core.Recognizers...)
}
