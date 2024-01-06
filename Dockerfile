# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.21.5 as builder

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
COPY sri/ sri/
COPY srictl/ srictl/
COPY srictl-bucket/ srictl-bucket/
COPY srictl-machine/ srictl-machine/
COPY srictl-volume/ srictl-volume/
COPY poollet/ poollet/
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

FROM builder as machinepoollet-builder

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -ldflags="-s -w" -a -o bin/machinepoollet ./poollet/machinepoollet/cmd/machinepoollet/main.go

FROM builder as srictl-machine-builder

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -ldflags="-s -w" -a -o bin/srictl-machine ./srictl-machine/cmd/srictl-machine/main.go

FROM builder as volumepoollet-builder

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -ldflags="-s -w" -a -o bin/volumepoollet ./poollet/volumepoollet/cmd/volumepoollet/main.go

FROM builder as srictl-volume-builder

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -ldflags="-s -w" -a -o bin/srictl-volume ./srictl-volume/cmd/srictl-volume/main.go

FROM builder as bucketpoollet-builder

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -ldflags="-s -w" -a -o bin/bucketpoollet ./poollet/bucketpoollet/cmd/bucketpoollet/main.go

FROM builder as srictl-bucket-builder

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -ldflags="-s -w" -a -o bin/srictl-bucket ./srictl-bucket/cmd/srictl-bucket/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot as manager
WORKDIR /
COPY --from=manager-builder /workspace/bin/controller-manager .
USER 65532:65532

ENTRYPOINT ["/controller-manager"]

FROM gcr.io/distroless/static:nonroot as apiserver
WORKDIR /
COPY --from=apiserver-builder /workspace/bin/apiserver .
USER 65532:65532

ENTRYPOINT ["/apiserver"]

FROM gcr.io/distroless/static:nonroot as machinepoollet
WORKDIR /
COPY --from=machinepoollet-builder /workspace/bin/machinepoollet .
USER 65532:65532

ENTRYPOINT ["/machinepoollet"]

FROM debian:bullseye-slim as srictl-machine
WORKDIR /
COPY --from=srictl-machine-builder /workspace/bin/srictl-machine .
USER 65532:65532

FROM gcr.io/distroless/static:nonroot as volumepoollet
WORKDIR /
COPY --from=volumepoollet-builder /workspace/bin/volumepoollet .
USER 65532:65532

ENTRYPOINT ["/volumepoollet"]

FROM debian:bullseye-slim as srictl-volume
WORKDIR /
COPY --from=srictl-volume-builder /workspace/bin/srictl-volume .
USER 65532:65532

FROM gcr.io/distroless/static:nonroot as bucketpoollet
WORKDIR /
COPY --from=bucketpoollet-builder /workspace/bin/bucketpoollet .
USER 65532:65532

ENTRYPOINT ["/bucketpoollet"]

FROM debian:bullseye-slim as srictl-bucket
WORKDIR /
COPY --from=srictl-bucket-builder /workspace/bin/srictl-bucket .
USER 65532:65532
