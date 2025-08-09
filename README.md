# Go Dumper

A self-hosted MySQL/MariaDB backup tool similar to MySQLDumper, built with Go and Vue.js.

## Features

- 🗄️ **Native Backup/Restore** - No external mysqldump dependency
- 🌐 **Web Interface** - Modern Vue.js frontend with TypeScript
- 📅 **Automated Scheduling** - Daily backups with customizable retention
- 🔐 **Encrypted Storage** - Secure password encryption with AES-GCM
- 🐳 **Docker Ready** - Multi-stage builds with multi-arch support
- ⚡ **High Performance** - Streaming backups with batch processing
- 🧪 **Well Tested** - Comprehensive unit, integration, and E2E tests

## Quick Start

### Docker (Recommended)

```bash
# Generate encryption key
export APP_ENC_KEY=$(openssl rand -base64 32)

# Run with Docker Compose
curl -o docker-compose.yml https://raw.githubusercontent.com/user/go-dumper/main/docker-compose.yml
docker compose up -d
```

### Docker Run

```bash
docker run -p 8080:8080 \
  -e APP_ENC_KEY=YOUR_32_BYTE_BASE64_KEY \
  -e BACKUP_DIR=/data/backups \
  -e SQLITE_PATH=/data/app/app.db \
  -v go-dumper-backups:/data/backups \
  -v go-dumper-data:/data/app \
  ghcr.io/user/go-dumper:latest
```

### From Source

```bash
git clone https://github.com/user/go-dumper.git
cd go-dumper

# Option 1: Full setup (requires Go + Node.js)
make setup-dev

# Option 2: Frontend-only setup (if Go is not available yet)
make setup-ui

# Generate encryption key and set up environment
make generate-key  # Copy this key
cp .env.local.example .env.local  # For local development
# Edit .env.local and paste the generated APP_ENC_KEY

# Start development servers
make dev
```

#### Prerequisites

- **Go 1.22+** - [Install from golang.org](https://golang.org/doc/install)
- **Node.js 20+** - [Install from nodejs.org](https://nodejs.org/)
- **Docker** (optional) - For testing and deployment

## Configuration

### Environment Variables

The application supports multiple ways to set environment variables:

1. **System environment variables** (highest priority)
2. **`.env.local` file** (local development overrides)
3. **`.env` file** (default configuration)

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | Server port | `8080` |
| `APP_ENC_KEY` | 32-byte base64 encryption key | **Required** |
| `SQLITE_PATH` | SQLite database path | `/data/app/app.db` |
| `BACKUP_DIR` | Backup storage directory | `/data/backups` |
| `ADMIN_USER` | Basic auth username (optional) | - |
| `ADMIN_PASS` | Basic auth password (optional) | - |

### Environment File Setup

```bash
# For production
cp .env.example .env

# For local development (overrides .env)
cp .env.local.example .env.local

# Generate encryption key
make generate-key
# or
openssl rand -base64 32
```

### Development vs Production

- **Development**: Use `.env.local` for local overrides (auto-loaded, git-ignored)
- **Production**: Use `.env` or system environment variables
- **Docker**: Environment variables passed via docker run or docker-compose

## Usage

1. **Access the web interface** at http://localhost:8080
2. **Add a target** - Configure your MySQL/MariaDB connection
3. **Create backups** - Manual or scheduled
4. **Download/Restore** - Manage your backup files

### API Endpoints

The tool provides a REST API for automation:

```bash
# List targets
curl http://localhost:8080/api/targets

# Create backup
curl -X POST http://localhost:8080/api/targets/1/backup

# Download backup
curl -o backup.sql.gz http://localhost:8080/api/backups/1/download

# Health check
curl http://localhost:8080/healthz
```

## Development

### Prerequisites

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose
- Make (optional)

### Setup

```bash
make setup-dev
make generate-key  # Add to .env file
make dev
```

### Testing

```bash
# All tests
make ci

# Unit tests only
make test

# Integration tests (requires Docker)
make test-int

# Frontend tests
make test-ui

# E2E tests
make test-e2e
```

### Building

```bash
# Local build
make build

# Docker build
make docker-build

# Docker run
make docker-run
```

## Architecture

### Backend (Go)
- **Gin** - HTTP framework
- **SQLite** - Application database (modernc.org/sqlite, CGO-free)
- **MySQL Driver** - github.com/go-sql-driver/mysql
- **Native Dumper** - Custom implementation with consistent snapshots

### Frontend (Vue.js)
- **Vue 3** - Progressive framework
- **TypeScript** - Type safety
- **Pinia** - State management
- **Vue Router** - Routing
- **Tailwind CSS** - Styling
- **DaisyUI** - Components

### Database Schema

#### Targets Table
Stores backup target configurations with encrypted passwords.

#### Backups Table
Tracks backup history, status, and file metadata with automatic cleanup.

### Backup Process

1. **Consistent Snapshot** - `REPEATABLE READ` isolation
2. **Schema Export** - `SHOW CREATE TABLE` for all tables
3. **Data Export** - Streaming with configurable batching
4. **Compression** - Optional gzip compression
5. **Cleanup** - Automatic rotation based on retention policy

## Security

- 🔐 **Encrypted Passwords** - AES-GCM encryption for database credentials
- 👤 **Optional Authentication** - Basic auth protection
- 🛡️ **SQL Injection Protection** - Parameterized queries
- 📝 **Security Scanning** - Trivy vulnerability scanning in CI

## Deployment

### Docker Compose (Production)

```yaml
version: '3.8'
services:
  go-dumper:
    image: ghcr.io/user/go-dumper:latest
    ports:
      - "8080:8080"
    environment:
      APP_ENC_KEY: YOUR_SECURE_KEY_HERE
      ADMIN_USER: admin
      ADMIN_PASS: secure_password
    volumes:
      - backups:/data/backups
      - data:/data/app
    restart: unless-stopped

volumes:
  backups:
  data:
```

### Kubernetes

See [examples/kubernetes/](examples/kubernetes/) for sample manifests.

### Reverse Proxy

```nginx
location / {
    proxy_pass http://localhost:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make ci`)
4. Commit changes (`git commit -m 'Add amazing feature'`)
5. Push to branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Troubleshooting

### Common Issues

#### Setup Issues

**Problem**: `make setup-dev` fails with "go: command not found"
```bash
# Solution: Install Go first, then run
make setup-ui    # Set up frontend first
# Install Go from https://golang.org/doc/install
make setup-dev   # Then complete setup
```

**Problem**: `npm ci` fails with missing package-lock.json
```bash
# Solution: Use npm install instead
cd web/app && npm install
# This generates package-lock.json automatically
```

**Problem**: DaisyUI or TailwindCSS version conflicts
```bash
# Solution: Clear npm cache and reinstall
cd web/app
rm -rf node_modules package-lock.json
npm install
```

#### Runtime Issues

**Problem**: "APP_ENC_KEY environment variable is required"
```bash
# Solution: Generate and set encryption key
make generate-key
# Copy the generated key to .env.local
echo "APP_ENC_KEY=your_generated_key_here" >> .env.local
```

**Problem**: Permission denied when creating backup directories
```bash
# Solution: Check permissions and create directories
mkdir -p ./data ./backups
# Or adjust BACKUP_DIR and SQLITE_PATH in .env.local
```

**Problem**: MySQL connection fails
- Check MySQL/MariaDB server is running
- Verify connection credentials
- Test connection manually: `mysql -h host -u user -p database`

#### Development Issues

**Problem**: Frontend dev server not starting
```bash
# Solution: Ensure Node.js and npm are installed
node --version  # Should be 20+
npm --version
cd web/app && npm run dev
```

**Problem**: Backend compilation errors
```bash
# Solution: Check Go version and dependencies
go version      # Should be 1.22+
go mod download
go mod tidy
```

### Getting Help

- 📖 **Documentation** - [GitHub Wiki](https://github.com/user/go-dumper/wiki)
- 🐛 **Issues** - [GitHub Issues](https://github.com/user/go-dumper/issues)
- 💬 **Discussions** - [GitHub Discussions](https://github.com/user/go-dumper/discussions)

## Acknowledgments

- Inspired by [MySQLDumper](http://www.mysqldumper.de/)
- Built with modern tools and practices
- Community-driven development