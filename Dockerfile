# Stage 1: Frontend bauen
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY web/app/package*.json ./
RUN npm ci
COPY web/app/ ./
RUN npm run build

# Stage 2: Backend bauen
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
RUN apk add --no-cache git

# 1) Mod-Dateien und Source rein
COPY go.mod go.sum ./
COPY internal/ ./internal/
COPY ["cmd/", "./cmd/"]

# 2) Module aufräumen + laden (fügt fehlende Einträge wie github.com/ncruces/go-strftime hinzu)
RUN go mod tidy && go mod download

# 3) Bauen
# (optional: explizit für Buildx)
# ARG TARGETOS TARGETARCH
# ENV GOOS=$TARGETOS GOARCH=$TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/app

# Stage 3: Finales Image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
RUN mkdir -p /data/app /data/backups
COPY --from=backend-builder /app/main .
COPY --from=frontend-builder /app/dist ./web/public
RUN addgroup -g 1001 -S appuser && \
    adduser -S -D -H -u 1001 -h /data -s /sbin/nologin -G appuser -g appuser appuser && \
    chown -R appuser:appuser /data
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1
CMD ["./main"]
