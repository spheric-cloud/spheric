// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package event

type CreateEvent[Object any] struct {
	Object Object
}

type UpdateEvent[Object any] struct {
	ObjectOld Object

	ObjectNew Object
}

type DeleteEvent[Object any] struct {
	Object Object
}

type GenericEvent[Object any] struct {
	Object Object
}
