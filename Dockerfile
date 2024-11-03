FROM golang:1.23-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build

FROM gcr.io/distroless/static AS bot

ENV GIN_MODE=release

WORKDIR /elf
COPY --from=builder /build/elf /elf/elf

ENTRYPOINT [ "/elf/elf", "run" ]
