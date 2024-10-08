// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package instanceevent

type Listener interface {
	Enqueue()
}

type EnqueueFunc struct {
	EnqueueFunc func()
}

func (n EnqueueFunc) Enqueue() {
	if n.EnqueueFunc != nil {
		n.EnqueueFunc()
	}
}

type ListenerRegistration interface{}
