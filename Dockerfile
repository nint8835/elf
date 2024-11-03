FROM golang:1.23-bookworm AS builder

WORKDIR /build
COPY . .
RUN go mod download && \
    go mod verify && \
    go build -o elf .

FROM debian:bookworm-slim

WORKDIR /elf
COPY --from=builder /build/elf /elf/elf
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "/elf/elf", "run" ]
