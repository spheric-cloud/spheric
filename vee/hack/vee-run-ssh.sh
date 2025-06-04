#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

SSH_HOST="${1:-remote-host}"
GOOS="${2:-linux}"
GOARCH="${3-amd64}"

cd "$SCRIPT_DIR"/../..

LOCAL_DIR="$HOME/.vee-run-ssh"
LOCAL_SOCKET="$LOCAL_DIR/vee.sock"
if [ -f "$LOCAL_SOCKET" ]; then
  echo "$LOCAL_SOCKET is already in-use"
  exit 1
fi

echo "Building vee binary"
BINARY="vee_${GOOS}_${GOARCH}"
GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=0 go build -o ./bin/"$BINARY" ./vee

echo "Preparing local socket directory"
mkdir -p "$LOCAL_DIR"

echo "Running vee on host"
trap 'rm -f $LOCAL_SOCKET' EXIT INT TERM
ssh "$SSH_HOST" 'pkill -f /opt/vee/bin/vee' || [ $? -eq 1 ]
scp -q bin/"$BINARY" "$SSH_HOST":/opt/vee/bin/vee
ssh -L "$LOCAL_SOCKET":/run/vee/vee.sock "$SSH_HOST" '/opt/vee/bin/vee'
