// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package runtime

import "iter"

type Object interface {
}

// List is a list of objects.
type List[Object any] interface {
	// Item returns a reference to the Object at the given index.
	Item(idx int) Object
	// All returns an iter.Seq over all Object references.
	All() iter.Seq[Object]
	// Len returns the length of the list.
	Len() int
}
