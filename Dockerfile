# --- Stage 1: Builder ---
FROM golang:1.26.1-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app

# Create the non-root user (matching the companion for consistency)
RUN adduser -D -u 1000 wolf

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the static binary
# We still build a static binary so it's portable and fast
# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gatus-bridge ./cmd/bridge/main.go

# --- Stage 2: Final Image ---
FROM alpine:3.23

# Install runtime essentials
# We keep alpine here to have 'sh' and 'chown' available for volume debugging
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

#  Copy user definitions
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy the binary
COPY --from=builder /app/gatus-bridge .

# Setup the data directory where Gatus config lives
# We ensure the wolf user owns this so it can write the YAML
RUN mkdir -p /data && chown -R wolf:wolf /data

# Standard environment variables
ENV GATUS_CONFIG_PATH=/data/config.yaml
ENV PORT=8080
ENV TZ=America/Montreal

USER wolf

# Healthcheck using wget (standard in Alpine)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

EXPOSE 8080

ENTRYPOINT ["./gatus-bridge"]