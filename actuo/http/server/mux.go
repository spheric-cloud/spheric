// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	MethodPost   = http.MethodPost
	MethodGet    = http.MethodGet
	MethodPut    = http.MethodPut
	MethodPatch  = http.MethodPatch
	MethodDelete = http.MethodDelete
	MethodWatch  = "WATCH"
)

var knownMethods = map[string]struct{}{
	MethodPost:   {},
	MethodGet:    {},
	MethodPut:    {},
	MethodPatch:  {},
	MethodDelete: {},
	MethodWatch:  {},
}

type getAndWatchHandler struct {
	getHandler   http.Handler
	watchHandler http.Handler
}

func (h *getAndWatchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	watch, _ := strconv.ParseBool(req.URL.Query().Get("watch"))
	if watch {
		if h.watchHandler == nil {
			http.NotFound(w, req)
			return
		}

		h.watchHandler.ServeHTTP(w, req)
		return
	}

	if h.getHandler == nil {
		http.NotFound(w, req)
		return
	}

	h.getHandler.ServeHTTP(w, req)
}

type ServeMux struct {
	mux                             *http.ServeMux
	getAndWatchHandlerByPathPattern map[string]*getAndWatchHandler
}

func NewServeMux() *ServeMux {
	mux := http.NewServeMux()
	return &ServeMux{
		mux:                             mux,
		getAndWatchHandlerByPathPattern: make(map[string]*getAndWatchHandler),
	}
}

func (m *ServeMux) Handle(pattern string, handler http.Handler) {
	if pattern == "" {
		panic("Invalid pattern")
	}
	if pattern[0] != '/' {
		method, pathPattern, found := strings.Cut(pattern, " ")
		if !found {
			panic("Invalid pattern")
		}

		if _, ok := knownMethods[method]; !ok {
			panic("Invalid method")
		}

		isWatch := strings.EqualFold(method, "watch")
		if isWatch {
			// Rewrite to regular http GET
			pattern = fmt.Sprintf("%s %s", http.MethodGet, pathPattern)
		}

		if isWatch || strings.EqualFold(method, MethodGet) {
			var getPtr func(*getAndWatchHandler) *http.Handler
			if isWatch {
				getPtr = func(g *getAndWatchHandler) *http.Handler { return &g.watchHandler }
			} else {
				getPtr = func(g *getAndWatchHandler) *http.Handler { return &g.getHandler }
			}

			if gwHandler, ok := m.getAndWatchHandlerByPathPattern[pathPattern]; ok {
				ptr := getPtr(gwHandler)
				if *ptr != nil {
					panic(fmt.Sprintf("Duplicate %s handler for pattern %s", method, pathPattern))
				}

				*ptr = handler
				return
			}

			gwHandler := &getAndWatchHandler{}
			*getPtr(gwHandler) = handler
			m.getAndWatchHandlerByPathPattern[pathPattern] = gwHandler
			handler = gwHandler
		}
	}

	m.mux.Handle(pattern, handler)
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.mux.ServeHTTP(w, req)
}
