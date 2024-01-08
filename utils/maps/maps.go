// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package maps

// Pop gets the value associated with the key (if any) and deletes it from the map.
func Pop[M ~map[K]V, K comparable, V any](m M, key K) (V, bool) {
	v, ok := m[key]
	delete(m, key)
	return v, ok
}
