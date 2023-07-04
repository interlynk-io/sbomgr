FROM golang:1.20-alpine AS builder
LABEL org.opencontainers.image.source="https://github.com/interlynk-io/sbomgr"

RUN apk add --no-cache make
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make ; make build

FROM scratch
LABEL org.opencontainers.image.source="https://github.com/interlynk-io/sbomgr"
LABEL org.opencontainers.image.description="SBOM Grep - Search through SBOMs"
LABEL org.opencontainers.image.licenses=Apache-2.0

COPY --from=builder /app/build/sbomgr /app/sbomgr
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "/app/sbomgr" ]