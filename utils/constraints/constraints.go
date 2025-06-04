// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package constraints

type Channel[T any] interface {
	~chan T | ~<-chan T | ~chan<- T
}

type SendChannel[T any] interface {
	~chan T | ~chan<- T
}

type ReceiveChannel[T any] interface {
	~chan T | ~<-chan T
}
