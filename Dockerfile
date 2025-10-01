FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o generator ./cmd/generator.go

# ===================================

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/generator .
COPY --from=builder /build/config ./config

EXPOSE 8000

CMD ["./generator", "-c", "config/config.yaml"]
