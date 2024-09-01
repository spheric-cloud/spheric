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

func Single[M ~map[K]V, K comparable, V any](m M) (K, V, bool) {
	if len(m) != 1 {
		var (
			zeroK K
			zeroV V
		)
		return zeroK, zeroV, false
	}
	for k, v := range m {
		return k, v, true
	}
	panic("Single: unreachable")
}

func MustSingle[M ~map[K]V, K comparable, V any](m M) (K, V) {
	k, v, ok := Single(m)
	if !ok {
		panic("MustSingle: not a single value")
	}
	return k, v
}
