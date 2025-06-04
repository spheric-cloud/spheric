// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package os

import (
	"fmt"
	"os"

	utilerrors "spheric.cloud/spheric/utils/errors"
)

func RemoveSocket(socketname string) error {
	stat, err := os.Stat(socketname)
	if err != nil {
		return err
	}
	if stat.Mode().Type()&os.ModeSocket == 0 {
		return fmt.Errorf("%s is not a socket", socketname)
	}

	return os.Remove(socketname)
}

func IgnoreNotExist(err error) error {
	return utilerrors.Ignore(err, os.ErrNotExist)
}

func EnsureSocketGone(socketname string) error {
	return IgnoreNotExist(RemoveSocket(socketname))
}
