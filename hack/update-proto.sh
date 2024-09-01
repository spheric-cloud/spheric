#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
REPO_ROOT="$SCRIPT_DIR/.."
export TERM="xterm-256color"

blue="$(tput setaf 4)"
normal="$(tput sgr0)"

PROTOC_GEN_GO="$PROTOC_GEN_GO"
PROTOC_GEN_GO_GRPC="$PROTOC_GEN_GO_GRPC"

function generate() {
  package="$1"
  (
  export PATH="$PATH:$(dirname "$PROTOC_GEN_GO")"
  export PATH="$PATH:$(dirname "$PROTOC_GEN_GO_GRPC")"
  echo "Generating ${blue}$package${normal}"
  protoc \
    --proto_path "." \
    --go_out="." \
    --go-grpc_out="." \
    --go_opt=module=spheric.cloud/spheric \
    --go-grpc_opt=module=spheric.cloud/spheric \
    "$package/api.proto"
  )
}

generate "iri-api/apis/runtime/v1alpha1"
