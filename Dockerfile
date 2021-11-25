FROM golang:1.17-bullseye AS builder

WORKDIR /build
COPY . .
RUN go mod download && \
    go mod verify && \
    go build -o elf ./cmd/elf

FROM debian:bullseye-slim

WORKDIR /elf
COPY --from=builder /build/elf /elf/elf
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "/elf/elf" ]
