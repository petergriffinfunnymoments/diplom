# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments
ARG SERVICE_NAME
ARG VERSION=dev
ARG BUILD_TIME

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" -o /app/bin/${SERVICE_NAME} ./cmd/${SERVICE_NAME}

# Final stage
FROM alpine:3.19 AS final

WORKDIR /app

# Install CA certificates and timezone
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
ARG SERVICE_NAME
COPY --from=builder /app/bin/${SERVICE_NAME} /app/${SERVICE_NAME}
COPY --from=builder /app/configs/ /app/configs/

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose port (default 8080)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the service
ENTRYPOINT ["/app/"]
CMD []

# ============================================
# Multi-stage builds for each service
# ============================================

# API Gateway
FROM builder AS api-gateway-builder
ARG SERVICE_NAME=api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/api-gateway ./cmd/api-gateway

FROM final AS api-gateway
ARG SERVICE_NAME=api-gateway
COPY --from=api-gateway-builder /app/bin/api-gateway /app/api-gateway
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/api-gateway"]

# Orchestrator
FROM builder AS orchestrator-builder
ARG SERVICE_NAME=orchestrator
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/orchestrator ./cmd/orchestrator

FROM final AS orchestrator
ARG SERVICE_NAME=orchestrator
COPY --from=orchestrator-builder /app/bin/orchestrator /app/orchestrator
USER appuser
ENTRYPOINT ["/app/orchestrator"]

# Validator
FROM builder AS validator-builder
ARG SERVICE_NAME=validator
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/validator ./cmd/validator

FROM final AS validator
ARG SERVICE_NAME=validator
COPY --from=validator-builder /app/bin/validator /app/validator
USER appuser
ENTRYPOINT ["/app/validator"]

# Anti-Fraud
FROM builder AS antifraud-builder
ARG SERVICE_NAME=antifraud
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/antifraud ./cmd/antifraud

FROM final AS antifraud
ARG SERVICE_NAME=antifraud
COPY --from=antifraud-builder /app/bin/antifraud /app/antifraud
USER appuser
ENTRYPOINT ["/app/antifraud"]

# Tokenizer
FROM builder AS tokenizer-builder
ARG SERVICE_NAME=tokenizer
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/tokenizer ./cmd/tokenizer

FROM final AS tokenizer
ARG SERVICE_NAME=tokenizer
COPY --from=tokenizer-builder /app/bin/tokenizer /app/tokenizer
USER appuser
ENTRYPOINT ["/app/tokenizer"]

# Adapter
FROM builder AS adapter-builder
ARG SERVICE_NAME=adapter
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/adapter ./cmd/adapter

FROM final AS adapter
ARG SERVICE_NAME=adapter
COPY --from=adapter-builder /app/bin/adapter /app/adapter
USER appuser
ENTRYPOINT ["/app/adapter"]

# Notification
FROM builder AS notification-builder
ARG SERVICE_NAME=notification
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/notification ./cmd/notification

FROM final AS notification
ARG SERVICE_NAME=notification
COPY --from=notification-builder /app/bin/notification /app/notification
USER appuser
ENTRYPOINT ["/app/notification"]

# Logger
FROM builder AS logger-builder
ARG SERVICE_NAME=logger
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/logger ./cmd/logger

FROM final AS logger
ARG SERVICE_NAME=logger
COPY --from=logger-builder /app/bin/logger /app/logger
USER appuser
ENTRYPOINT ["/app/logger"]
