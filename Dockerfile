# Build stage
FROM --platform=$BUILDPLATFORM golang:1.27.4@sha256:205c7636d74e0ffbe3d007162df7d74f859f5a8c8322dfb9f907977686c126cf AS builder

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev

WORKDIR /workspace

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" \
    go build -o /workspace/provider \
    -ldflags "-s -w -X main.Version=${VERSION}" \
    ./cmd/provider

# Download OpenTofu binary for the target architecture
ARG OPENTOFU_VERSION=1.9.1
RUN apt-get update && apt-get install -y --no-install-recommends unzip && \
    curl -fsSL "https://github.com/opentofu/opentofu/releases/download/v${OPENTOFU_VERSION}/tofu_${OPENTOFU_VERSION}_${TARGETOS}_${TARGETARCH}.zip" \
    -o /tmp/tofu.zip && \
    unzip /tmp/tofu.zip -d /tmp/tofu && \
    mv /tmp/tofu/tofu /workspace/tofu && \
    chmod +x /workspace/tofu && \
    rm -rf /tmp/tofu /tmp/tofu.zip

# Runtime stage
FROM gcr.io/distroless/static:nonroot@sha256:e3f945647ffb95b5839c07038d64f9811adf17308b9121d8a2b87b6a22a80a39

COPY --from=builder /workspace/provider /usr/local/bin/provider
COPY --from=builder /workspace/tofu /usr/local/bin/terraform

USER 65532:65532

ENTRYPOINT ["provider"]
