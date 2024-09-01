#!/usr/bin/env bash

K8S_VERSION="${1:-"v0.31.0"}"

go get \
  "k8s.io/api@$K8S_VERSION" \
  "k8s.io/apimachinery@$K8S_VERSION" \
  "k8s.io/apiserver@$K8S_VERSION" \
  "k8s.io/client-go@$K8S_VERSION" \
  "k8s.io/component-base@$K8S_VERSION" \
  "k8s.io/kube-aggregator@$K8S_VERSION" \
  "k8s.io/kubectl@$K8S_VERSION"
