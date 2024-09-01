// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"io"
	"os"
)

func CopyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		_ = dstFile.Close()
		return err
	}
	if err := dstFile.Close(); err != nil {
		return err
	}
	return nil
}

func ToMap[S ~[]E, E any, K comparable, V any](s S, f func(E) (K, V)) map[K]V {
	res := make(map[K]V)
	for _, e := range s {
		k, v := f(e)
		res[k] = v
	}
	return res
}

func ToMapByKey[S ~[]E, E any, K comparable](s S, f func(E) K) map[K]E {
	return ToMap[S, E, K, E](s, func(e E) (K, E) {
		return f(e), e
	})
}
