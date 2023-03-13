# Build the manager binary
FROM golang:1.19 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor/ vendor/
COPY types/ types/
COPY storage/ storage/
COPY internal/ internal/
COPY config/ config/
COPY cmd/ cmd/
COPY api/ api/

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o=bin/sportsnews -mod=vendor cmd/sportsnews/main.go

# copy licenses
RUN mkdir /licenses
COPY LICENSE /LICENSE

# have default storageKind set to memory
RUN echo 'storageKind: "memory"' > /workspace/config/custom.yaml

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/bin/sportsnews /
COPY --from=builder /workspace/config /config
COPY --from=builder /licenses /licenses
EXPOSE 8080
USER 65532:65532

ENV USER_UID=65532

ENTRYPOINT ["/sportsnews"]