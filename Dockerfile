# Use buildx for multi-platform builds
# Build stage
FROM --platform=$BUILDPLATFORM golang:1.22.2-alpine AS builder
LABEL org.opencontainers.image.source="https://github.com/interlynk-io/sbomgr"

RUN apk add --no-cache make git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build for multiple architectures
ARG TARGETOS TARGETARCH
RUN make build && chmod +x ./build/sbomgr

# Final stage
FROM alpine:3.19
LABEL org.opencontainers.image.source="https://github.com/interlynk-io/sbomgr"
LABEL org.opencontainers.image.description="Search through SBOMs"
LABEL org.opencontainers.image.licenses=Apache-2.0

# Copy our static executable
COPY --from=builder /app/build/sbomgr /app/sbomgr

# Disable version check
ENV INTERLYNK_DISABLE_VERSION_CHECK=true

ENTRYPOINT [ "/app/sbomgr" ]
