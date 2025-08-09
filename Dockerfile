# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# Copy package files
COPY web/app/package*.json ./
RUN npm ci --only=production

# Copy source code
COPY web/app/ ./

# Build frontend
RUN npm run build

# Stage 2: Build backend
FROM golang:1.22-alpine AS backend-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY cmd/ cmd/
COPY internal/ internal/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Stage 3: Final image
FROM alpine:latest

# Install ca-certificates for SSL/TLS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Create directories for data
RUN mkdir -p /data/app /data/backups

# Copy the binary from builder stage
COPY --from=backend-builder /app/main .

# Copy the built frontend from frontend-builder stage
COPY --from=frontend-builder /app/dist ./web/public

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -S -D -H -u 1001 -h /data -s /sbin/nologin -G appuser -g appuser appuser

# Change ownership of data directories
RUN chown -R appuser:appuser /data

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

# Command to run
CMD ["./main"]