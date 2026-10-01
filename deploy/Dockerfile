FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/http \
    ./cmd/http

FROM alpine:3.21

RUN apg add --no-cahce ca-certificates \
    && adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/bin/http /app/http

USER appuser
EXPOSE 5555

CMD ["/app/http"]
