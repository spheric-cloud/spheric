// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"net/http"
	"net/url"
	"strconv"

	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/util/proxy"
	ctrl "sigs.k8s.io/controller-runtime"
	sri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	"spheric.cloud/spheric/spherelet/api/v1alpha1"
)

func (s *Server) serveExec(w http.ResponseWriter, req *http.Request, namespace, name string) {
	ctx := req.Context()
	log := ctrl.LoggerFrom(ctx)

	listInstancesRes, err := s.runtimeService.ListInstances(ctx, &sri.ListInstancesRequest{
		Filter: &sri.InstanceFilter{
			LabelSelector: map[string]string{
				v1alpha1.InstanceNamespaceLabel: namespace,
				v1alpha1.InstanceNameLabel:      name,
			},
		},
	})
	if err != nil {
		log.Error(err, "Error listing instances")
		s.writeError(w, err)
		return
	}
	if len(listInstancesRes.Instances) == 0 {
		http.Error(w, "instance not found", http.StatusNotFound)
		return
	}

	instance := listInstancesRes.Instances[0]
	execRes, err := s.runtimeService.Exec(ctx, &sri.ExecRequest{
		InstanceId: instance.Metadata.Id,
	})
	if err != nil {
		log.Error(err, "Error getting exec url")
		s.writeError(w, err)
		return
	}

	execURL, err := url.Parse(execRes.Url)
	if err != nil {
		log.Error(err, "Error parsing exec url")
		s.writeError(w, err)
		return
	}

	proxyStream(w, req, execURL)
}

func (s *Server) writeError(w http.ResponseWriter, err error) {
	status, _ := grpcstatus.FromError(err)
	var code int
	switch status.Code() {
	case codes.NotFound:
		code = http.StatusNotFound
	case codes.ResourceExhausted:
		w.Header().Set("Retry-After", strconv.Itoa(int(s.cacheTTL.Seconds())))
	default:
		code = http.StatusInternalServerError
	}
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}

func proxyStream(w http.ResponseWriter, req *http.Request, url *url.URL) {
	handler := proxy.NewUpgradeAwareHandler(url, nil, false, true, &responder{})
	handler.ServeHTTP(w, req)
}

type responder struct{}

func (r *responder) Error(w http.ResponseWriter, req *http.Request, err error) {
	ctx := req.Context()
	log := ctrl.LoggerFrom(ctx)
	log.Error(err, "Error while proxying request")
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
