#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

err() { echo "$@" 1>&2; }

if ! which limactl > /dev/null; then
  err "limactl not installed"
  exit 1
fi

limactl start \
  --name vee \
  --tty=false \
  "$SCRIPT_DIR/veem.yaml"
