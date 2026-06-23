# Build stage: compile kube-agents binary.
ARG GO_VERSION=1.26.3-alpine3.23
FROM golang:${GO_VERSION} AS builder

# Set Go Proxy for faster dependency downloads.
ARG GOPROXY="https://goproxy.cn,direct"
ENV GOPROXY=${GOPROXY}

WORKDIR /workspace

# Copy go.mod and go.sum first for better layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the project source.
COPY . .

# Build arguments injected at build time.
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG GITCOMMIT
ARG BUILDDATE

# Build the binary with stripped symbols and version info.
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags "-s -w \
    -X github.com/codefuture-io/kube-agents/version.Version=${VERSION} \
    -X github.com/codefuture-io/kube-agents/version.GitCommit=${GITCOMMIT} \
    -X github.com/codefuture-io/kube-agents/version.BuildDate=${BUILDDATE}" \
    -a -o kube-agents ./cmd/kube-agents/

# Runtime stage: minimal distroless image.
FROM gcr.io/distroless/static-debian13:latest

WORKDIR /

COPY --from=builder /workspace/kube-agents /usr/local/bin/

# Copy config template and skills.
COPY --from=builder /workspace/config/ /etc/kube-agents/config/
COPY --from=builder /workspace/skills/ /skills/

# Run as non-root user.
USER 65534:65534

ENTRYPOINT ["/usr/local/bin/kube-agents"]
CMD ["serve"]
