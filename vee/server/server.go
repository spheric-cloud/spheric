// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"spheric.cloud/spheric/actuo/builder"
	"spheric.cloud/spheric/actuo/http/server"
	storagestore "spheric.cloud/spheric/actuo/storage/store"
	"spheric.cloud/spheric/vee/api"
)

type Server struct {
	store storagestore.Store[string, *api.Instance]
}

func New(filename string, store storagestore.Store[string, *api.Instance]) (*server.Server, error) {
	handler, err := builder.NewNamespacedMeta[*api.Instance](
		"instances",
		store,
	).
		CRUD().
		Build()
	if err != nil {
		return nil, err
	}

	return server.NewServer(filename, handler), nil
}
