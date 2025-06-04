// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package errors

import "errors"

func Ignore(err error, toIgnore ...error) error {
	if err == nil {
		return nil
	}
	for _, e := range toIgnore {
		if errors.Is(err, e) {
			return nil
		}
	}
	return err
}
