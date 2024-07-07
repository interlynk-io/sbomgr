FROM golang:1.22.2-alpine AS builder
LABEL org.opencontainers.image.source="https://github.com/interlynk-io/sbomgr"

RUN apk add --no-cache make git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN make ; make build

FROM scratch
LABEL org.opencontainers.image.source="https://github.com/interlynk-io/sbomgr"
LABEL org.opencontainers.image.description="Search through SBOMs"
LABEL org.opencontainers.image.licenses=Apache-2.0

COPY --from=builder /bin/sh /bin/grep /bin/busybox /bin/touch /bin/chmod /bin/mkdir /bin/date /bin/cat /bin/
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /lib/ld-musl-x86_64.so.1 /lib/ld-musl-x86_64.so.1
COPY --from=builder /tmp /tmp
COPY --from=builder /usr/bin /usr/bin

# Copy our static executable
COPY --from=builder /app/build/sbomgr /app/sbomgr

# Disable version check
ENV INTERLYNK_DISABLE_VERSION_CHECK=true

ENTRYPOINT [ "/app/sbomgr" ]
