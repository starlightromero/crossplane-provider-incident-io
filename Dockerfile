# Build stage
FROM --platform=$BUILDPLATFORM golang:1.26.3@sha256:313faae491b410a35402c05d35e7518ae99103d957308e940e1ae2cfa0aac29b AS builder

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
FROM gcr.io/distroless/static:nonroot@sha256:1c2c046bc09ed40fad370b599a0b1ae7987f55b01e247cf27a7c27cd97e5bbc7

COPY --from=builder /workspace/provider /usr/local/bin/provider
COPY --from=builder /workspace/tofu /usr/local/bin/terraform

USER 65532:65532

ENTRYPOINT ["provider"]
