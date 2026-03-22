# ---------- BUILD STAGE ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git and bash (for go mod)
RUN apk add --no-cache git bash

ENV CGO_ENABLED=0

# Download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build optimized binary
RUN go build -ldflags="-s -w" -trimpath -o business-api

# Fix file permissions in build stage for .env and keys
RUN chmod 644 /app/.env && \
    chmod -R 755 /app/keys

# ---------- RUNTIME STAGE ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy built binary and prepared files
COPY --from=builder /app/business-api .
COPY --from=builder /app/.env .
COPY --from=builder /app/keys ./keys

# Use non-root user
USER nonroot:nonroot

# Expose API port
EXPOSE 8080

# Entrypoint
ENTRYPOINT ["/app/business-api"]