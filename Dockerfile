# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.23 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY api/ api/
COPY client-go/ client-go/
COPY cmd/ cmd/
COPY internal/ internal/
COPY iri-api/ iri-api/
COPY spherelet/ spherelet/
COPY utils/ utils/

ARG TARGETOS
ARG TARGETARCH

RUN mkdir bin

FROM builder as apiserver-builder

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -ldflags="-s -w" -a -o bin/apiserver ./cmd/apiserver

FROM builder as manager-builder

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -ldflags="-s -w" -a -o bin/controller-manager ./cmd/controller-manager

FROM builder as spherelet-builder

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -ldflags="-s -w" -a -o bin/spherelet ./spherelet/cmd/spherelet/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot as controller-manager
WORKDIR /
COPY --from=manager-builder /workspace/bin/controller-manager .
USER 65532:65532

ENTRYPOINT ["/controller-manager"]

FROM gcr.io/distroless/static:nonroot as apiserver
WORKDIR /
COPY --from=apiserver-builder /workspace/bin/apiserver .
USER 65532:65532

ENTRYPOINT ["/apiserver"]

FROM gcr.io/distroless/static:nonroot as spherelet
WORKDIR /
COPY --from=spherelet-builder /workspace/bin/spherelet .
USER 65532:65532

ENTRYPOINT ["/spherelet"]
