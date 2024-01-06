#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
REPO_ROOT="$SCRIPT_DIR/.."
export TERM="xterm-256color"

blue="$(tput setaf 4)"
normal="$(tput sgr0)"

VGOPATH="$VGOPATH"
PROTOC_GEN_GOGO="$PROTOC_GEN_GOGO"

VIRTUAL_GOPATH="$(mktemp -d)"
trap 'rm -rf "$VIRTUAL_GOPATH"' EXIT

# Setup virtual GOPATH so the codegen tools work as expected.
(
cd "$REPO_ROOT"
"$VGOPATH" -o "$VIRTUAL_GOPATH"
)

export GOPATH="$VIRTUAL_GOPATH"
export GO111MODULE=off

function generate() {
  package="$1"
  (
  cd "$VIRTUAL_GOPATH/src"
  export PATH="$PATH:$(dirname "$PROTOC_GEN_GOGO")"
  echo "Generating ${blue}$package${normal}"
  protoc \
    --proto_path "./spheric.cloud/spheric/$package" \
    --proto_path "$VIRTUAL_GOPATH/src" \
    --gogo_out=plugins=grpc:"$VIRTUAL_GOPATH/src" \
    "./spheric.cloud/spheric/$package/api.proto"
  )
}

generate "sri/apis/meta/v1alpha1"
generate "sri/apis/machine/v1alpha1"
generate "sri/apis/volume/v1alpha1"
generate "sri/apis/bucket/v1alpha1"
