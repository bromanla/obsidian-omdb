# Stage 1: Builder
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build \
  -a \
  -ldflags="-w -s" \
  -o omdb-bot \
  ./cmd

# Stage 2: Runtime
FROM gcr.io/distroless/static-debian12
WORKDIR /app

COPY --from=builder /app/omdb-bot .
CMD ["./omdb-bot"]