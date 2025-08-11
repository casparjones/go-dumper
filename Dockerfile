# Stage 1: Frontend bauen
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY web/app/package*.json ./
RUN npm ci                # devDeps nötig, damit vite verfügbar ist
COPY web/app/ ./
RUN npm run build         # erzeugt /app/dist

# Stage 2: Backend bauen
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/app

# Stage 3: Finales Image (Go liefert statische Files aus)
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Verzeichnisse
RUN mkdir -p /data/app /data/backups

# Binary kopieren
COPY --from=backend-builder /app/main .

# Statische Assets dorthin kopieren, wo dein Router sie erwartet
COPY --from=frontend-builder /app/dist ./web/public

# Non-root User
RUN addgroup -g 1001 -S appuser && \
    adduser -S -D -H -u 1001 -h /data -s /sbin/nologin -G appuser -g appuser appuser && \
    chown -R appuser:appuser /data
USER appuser

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

CMD ["./main"]
