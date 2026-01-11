# ---------- Build stage ----------
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/api ./cmd/api

# ---------- Runtime stage ----------
FROM alpine:3.20

# Certs for HTTP clients if needed
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /bin/api /app/api

EXPOSE 8080

# Healthcheck uses the existing endpoint
HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8080/healthz || exit 1

ENTRYPOINT ["/app/api"]
