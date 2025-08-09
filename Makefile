.PHONY: help build test test-int test-ui test-e2e ci clean dev docker-build docker-run format lint setup-dev setup-ui generate-key check-tools info

# Default target
help:
	@echo "Go Dumper - MySQL/MariaDB Backup Tool"
	@echo ""
	@echo "Available targets:"
	@echo "  help        Show this help message"
	@echo "  build       Build the application"
	@echo "  dev         Run in development mode"
	@echo "  setup-dev   Set up development environment"
	@echo "  setup-ui    Set up frontend only (no Go required)"
	@echo "  test        Run Go unit tests with coverage"
	@echo "  test-int    Run Go integration tests"
	@echo "  test-ui     Run frontend tests with coverage"
	@echo "  test-e2e    Run end-to-end tests"
	@echo "  ci          Run all tests (CI pipeline)"
	@echo "  format      Format code"
	@echo "  lint        Run linters"
	@echo "  clean       Clean build artifacts"
	@echo "  docker-build Build Docker image"
	@echo "  docker-run   Run with Docker Compose"
	@echo "  generate-key Generate encryption key"
	@echo "  check-tools  Check if required tools are installed"

# Build the application
build:
	@echo "Building frontend..."
	cd web/app && (test -f package-lock.json && npm ci || npm install) && npm run build
	@echo "Building backend..."
	CGO_ENABLED=0 go build -o bin/go-dumper ./cmd/app

# Development mode
dev:
	@echo "Starting development servers..."
	@echo "Backend will run on :8080, Frontend on :5173"
	@echo "Make sure to set required environment variables:"
	@echo "  export APP_ENC_KEY=\$$(openssl rand -base64 32)"
	@echo ""
	@echo "Starting backend..."
	go run ./cmd/app &
	@echo "Starting frontend..."
	cd web/app && npm run dev

# Run Go unit tests with coverage
test:
	@echo "Running Go unit tests..."
	go test ./... -race -coverprofile=coverage.out
	@echo "Checking coverage threshold..."
	go tool cover -func=coverage.out | awk '/total:/ { \
		split($$3,a,"%"); \
		if (a[1]+0 < 80) { \
			print "❌ Coverage too low:" $$3 " (minimum: 80%)"; \
			exit 1 \
		} else { \
			print "✅ Coverage:" $$3 \
		} \
	}'

# Run integration tests
test-int:
	@echo "Starting test databases..."
	docker compose -f test/docker-compose.yml up -d
	@echo "Waiting for databases to be ready..."
	sleep 10
	@echo "Running integration tests..."
	go test ./internal/... -tags=integration -v || (docker compose -f test/docker-compose.yml down -v && exit 1)
	@echo "Cleaning up test databases..."
	docker compose -f test/docker-compose.yml down -v

# Run frontend tests
test-ui:
	@echo "Running frontend tests..."
	cd web/app && (test -f package-lock.json && npm ci || npm install) && npm run test:coverage

# Run E2E tests
test-e2e:
	@echo "Installing Playwright browsers..."
	cd web/app && (test -f package-lock.json && npm ci || npm install) && npx playwright install --with-deps
	@echo "Running E2E tests..."
	cd web/app && npm run test:e2e

# Run all tests (CI pipeline)
ci: test test-int test-ui test-e2e
	@echo "✅ All tests passed!"

# Format code
format:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "Formatting frontend code..."
	cd web/app && npm run format || true

# Run linters
lint:
	@echo "Running Go linters..."
	go vet ./...
	@echo "Running frontend linters..."
	cd web/app && npm run lint || true

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf web/public/*
	rm -rf web/app/dist/
	rm -rf web/app/node_modules/
	rm -f coverage.out

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t go-dumper:latest .

# Run with Docker Compose
docker-run:
	@echo "Starting Go Dumper with Docker Compose..."
	@echo "Make sure to update the APP_ENC_KEY in docker-compose.yml"
	@echo "Generate a new key with: openssl rand -base64 32"
	docker compose up -d

# Generate encryption key
generate-key:
	@echo "Generated encryption key:"
	@openssl rand -base64 32

# Development setup
setup-dev: check-tools
	@echo "Setting up development environment..."
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Installing frontend dependencies..."
	cd web/app && npm install
	@echo "Creating local development directories..."
	@mkdir -p ./data ./backups
	@echo ""
	@echo "✅ Development environment set up!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Copy .env.local.example to .env.local for local development"
	@echo "   OR copy .env.example to .env for basic setup"
	@echo "2. Generate a key with: make generate-key"
	@echo "3. Edit your .env/.env.local file and set APP_ENC_KEY"
	@echo "4. Run with: make dev"
	@echo ""
	@echo "Environment file priority:"
	@echo "  System env > .env.local > .env > defaults"

# Setup frontend only (useful if Go is not available)
setup-ui:
	@echo "Setting up frontend environment..."
	@command -v node >/dev/null 2>&1 || { echo "❌ Node.js is not installed. Please install Node.js 20+ from https://nodejs.org/"; exit 1; }
	@command -v npm >/dev/null 2>&1 || { echo "❌ npm is not installed. Please install npm"; exit 1; }
	@echo "Installing frontend dependencies..."
	cd web/app && npm install
	@echo "Creating local development directories..."
	@mkdir -p ./data ./backups
	@echo ""
	@echo "✅ Frontend environment set up!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Install Go 1.22+ from https://golang.org/doc/install"
	@echo "2. Run: make setup-dev (to complete full setup)"
	@echo "3. Generate a key with: make generate-key"
	@echo "4. Copy .env.local.example to .env.local and set APP_ENC_KEY"

# Check if all required tools are installed
check-tools:
	@echo "Checking required tools..."
	@command -v go >/dev/null 2>&1 || { echo "❌ Go is not installed. Please install Go 1.22+ from https://golang.org/doc/install"; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "❌ Node.js is not installed. Please install Node.js 20+ from https://nodejs.org/"; exit 1; }
	@command -v npm >/dev/null 2>&1 || { echo "❌ npm is not installed. Please install npm"; exit 1; }
	@echo "✅ Required tools are installed"
	@echo "Go version: $$(go version)"
	@echo "Node version: $$(node --version)"
	@echo "npm version: $$(npm --version)"

# Show project info
info:
	@echo "Go Dumper - MySQL/MariaDB Backup Tool"
	@echo "======================================"
	@echo ""
	@echo "Project structure:"
	@echo "  cmd/app/           - Main application"
	@echo "  internal/          - Internal packages"
	@echo "    ├── backup/      - Backup/restore logic"
	@echo "    ├── http/        - HTTP handlers and routing"
	@echo "    ├── scheduler/   - Automated backup scheduling"
	@echo "    └── store/       - Database models and repository"
	@echo "  web/app/           - Vue.js frontend"
	@echo "  test/              - Integration test setup"
	@echo ""
	@echo "Key features:"
	@echo "  - Native MySQL/MariaDB backup (no mysqldump dependency)"
	@echo "  - Web-based management interface"
	@echo "  - Automated scheduling with retention policies"
	@echo "  - Encrypted password storage"
	@echo "  - Docker support with multi-arch builds"
	@echo "  - Comprehensive test coverage"