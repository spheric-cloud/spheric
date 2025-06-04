// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package proto

import (
	"google.golang.org/protobuf/proto"
)

func Clone[message proto.Message](msg message) message {
	return proto.Clone(msg).(message)
}
