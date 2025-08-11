Erzeuge ein lauffähiges Open-Source Web-Tool ähnlich „MySQLDumper“.

1) Ziel
Selbstgehosteter Service zum Sichern und Wiederherstellen von MySQL/MariaDB. Mehrere Ziele verwalten. Backups lokal in konfigurierbarem Ordner (Docker-Volume). Keine Cloud-Sync-Funktion.

2) Tech
Backend: Go 1.22+, Gin (JSON-API, kein Server-Side-Templating).
Frontend: Vue 3 + Vite + Vue Router + Pinia, TypeScript.
UI: Tailwind CSS v4 + DaisyUI v5, lokal gebaut ohne CDNs.
Persistenz (App-Settings): SQLite (modernc.org/sqlite, ohne CGO).
MySQL/MariaDB Treiber: github.com/go-sql-driver/mysql.
Container: Docker (Multi-Stage).
CI: GitHub Actions baut Frontend und Backend, pusht Image zu GHCR.
Lizenz: MIT.

3) Kernfeatures (wie MySQLDumper)
Ziele anlegen/bearbeiten/löschen: Name, Host, Port, DB-Name, User, Passwort, Kommentar, tägliche Uhrzeit (HH:MM, UTC), Aufbewahrungstage, Komprimierung.
Manuelles Backup.
Backups je Ziel: Liste, Download .sql.gz, Löschen, Details.
Restore mit Bestätigungsdialog.
Interner Daily-Scheduler pro Ziel zur Uhrzeit.
Rotation nach Aufbewahrungstagen.
Health-/Readiness: /healthz, /readyz.
Optional Basic-Auth über ADMIN_USER/ADMIN_PASS (gilt für API; UI kann bei gesetzter Auth nur mit Login arbeiten).

4) Dump/Restore (ohne externes mysqldump)
Konsistenter Dump: REPEATABLE READ, START TRANSACTION WITH CONSISTENT SNAPSHOT (falls verfügbar).
Schema: SHOW CREATE TABLE; nach Möglichkeit Views/Triggers.
Daten: tabellenweise streamend; INSERT-Batching; korrektes Escaping.
FKs temporär deaktivieren und am Ende aktivieren.
Ausgabe UTF-8 .sql → gzip zu .sql.gz.
Dateiname: {targetName}_{YYYY-MM-DD_HH-mm-ss}.sql.gz.
BACKUP_DIR Default /data/backups.
Restore: .sql.gz entpacken und sequenziell ausführen; UI-Warnung.

5) Datenmodell (SQLite)
targets(id, name, host, port, dbname, user, password_enc, comment, schedule_time, retention_days, auto_compress, created_at, updated_at).
backups(id, target_id, started_at, finished_at, size_bytes, status[success|failed|running], file_path, notes).
password_enc per AES-GCM; Schlüssel aus ENV APP_ENC_KEY (32-Byte Base64).

6) API-Routen (JSON, Prefix /api)
GET /api/targets; POST /api/targets; GET /api/targets/:id; PUT /api/targets/:id; DELETE /api/targets/:id.
POST /api/targets/:id/backup.
GET /api/targets/:id/backups.
GET /api/backups/:id/download.
POST /api/backups/:id/restore.
DELETE /api/backups/:id.
GET /healthz; GET /readyz.
Static UI unter / (Vue SPA, History-Mode), API unter /api/*.
Dev: Vite-Proxy auf /api → :8080; CORS nur im Dev nötig.

7) Frontend (Vue)
Seiten: Dashboard, Ziele (Liste, Neu, Bearbeiten), Ziel-Backups (Liste, Details, Aktionen), Einstellungen.
State mit Pinia; API-Client (Axios/Fetch) mit Interceptors (Basic-Auth, Fehler-Handling).
Formvalidierung; Toaster-Feedback; Loading-States.
Tailwind v4 + DaisyUI v5; Build per PostCSS. Keine externen CDNs.
Projektstruktur:
web/app/
  src/main.ts, App.vue
  src/router/index.ts
  src/stores/index.ts
  src/pages/{Dashboard,Targets,TargetEdit,Backups,Settings}.vue
  src/components/{Forms,Tables,Toasts,Modals}.vue
  src/services/api.ts
  index.html
Build-Output nach web/public (von Gin als static served).

8) Backend-Struktur
cmd/app/main.go
internal/http/router.go, handlers/*.go, middleware/*.go
internal/backup/ (dump, restore, gzip, rotation)
internal/scheduler/ (daily jobs)
internal/store/ (SQLite models, migrations, repo)
web/public/ (gebautes Frontend)

9) Konfiguration (ENV)
APP_PORT (Default 8080)
APP_ENC_KEY (32-Byte Base64)
SQLITE_PATH (Default /data/app/app.db)
BACKUP_DIR (Default /data/backups)
ADMIN_USER, ADMIN_PASS (optional)

10) Docker
Multi-Stage:
1) Node-Stage: Vite-Build für Vue + Tailwind/DaisyUI; Ausgabe nach web/public.
2) Go-Stage: CGO_ENABLED=0; statisches Binary.
3) Final: distroless oder scratch/alpine; enthält Binary + web/public.
User non-root; EXPOSE 8080.
Volumes: /data/app (SQLite), /data/backups (Dumps).

11) GitHub Actions (.github/workflows/docker.yml)
Trigger: push auf main, Tags v*.
Schritte: Checkout; Setup Node; Frontend build; Setup QEMU+Buildx; Go build (falls nötig); Login GHCR; Docker Buildx build & push für linux/amd64 und linux/arm64.
Tags:
ghcr.io/<OWNER>/<REPO>:latest
ghcr.io/<OWNER>/<REPO>:{shortSHA}
bei Git-Tag zusätzlich ghcr.io/<OWNER>/<REPO>:{gitTag}

12) Qualität
Kontext-Timeouts und sauberes Error-Handling.
Unit-Tests: Dump-Serialisierung, Rotation, einfache Handler.
Linting und go vet.
README mit Beispielbefehlen:
docker run -p 8080:8080 \
  -e BACKUP_DIR=/data/backups \
  -e SQLITE_PATH=/data/app/app.db \
  -e APP_ENC_KEY=BASE64_32_BYTES_KEY \
  -v backups:/data/backups \
  -v app:/data/app \
  ghcr.io/OWNER/REPO:latest
Beschreibung Dev-Setup: Backend :8080, Vite :5173 mit Proxy auf /api.

13) Tests & TDD

A) Teststrategie
1. Backend Unit-Tests: Dump-Serialisierung, Restore-Executor, Rotation, Scheduler (Zeitpunkte), AES-GCM Verschlüsselung/Entschlüsselung, Repositories (SQLite).
2. Backend Integrations-Tests: Gegen echte MySQL/MariaDB mit testcontainers-go oder docker-compose. Szenarien: Konsistenter Dump großer Tabellen, Restore inkl. FK-Handling, Abbruch bei Timeout.
3. Frontend Unit- und Komponententests: Vitest für Stores, Services (API-Client), Form-Validierung und kritische UI-Komponenten.
4. E2E-Tests: Playwright für die wichtigsten Flows (Target anlegen, manuelles Backup, Download, Restore-Guard, Rotation sichtbar).
5. Zielwerte: Backend Coverage ≥ 80%, Frontend Coverage ≥ 70%. Builds sollen bei Unterschreitung fehlschlagen.

B) TDD-Vorgehen
1. Vor jeder neuen Funktion zuerst Tests schreiben.
2. Minimalen Code implementieren, bis die Tests grün sind.
3. Refactoring mit weitergrünen Tests.
4. Erst danach nächstes To-do beginnen.

C) Lokale Befehle (Makefile)
- make test           → Go Unit-Tests mit Coverage
- make test-int       → Go Integrations-Tests gegen MySQL in Docker
- make test-ui        → Vitest mit Coverage
- make test-e2e       → Playwright E2E
- make ci             → test, test-int, test-ui, test-e2e in Reihenfolge

D) CI-Anforderungen
1. Tests müssen VOR dem Docker-Build laufen und grün sein.
2. GitHub Actions führt aus:
   - Go: go test ./... -race -coverprofile=coverage.out und Abbruch wenn Coverage < 80%.
   - Node: npm ci, npm run test:coverage und Abbruch wenn Coverage < 70%.
   - E2E: npm run test:e2e (Headless), vorher Playwright-Browser installieren.
3. Nur bei Erfolg: Docker Buildx Build & Push.

E) Beispiel Makefile-Ziele
test: go test ./... -race -coverprofile=coverage.out && go tool cover -func=coverage.out | awk '/total:/ { if ($$3+0 < 80) exit 1 }'
test-int: docker compose -f test/docker-compose.yml up -d && go test ./internal/... -tags=integration -v && docker compose -f test/docker-compose.yml down -v
test-ui: npm run test:coverage
test-e2e: npx playwright install --with-deps && npm run test:e2e
ci: make test && make test-int && make test-ui && make test-e2e

F) GitHub Actions Snippet (Auszug)
- name: Go tests
  run: |
    go test ./... -race -coverprofile=coverage.out
- name: Node tests
  run: |
    npm ci
    npm run test:coverage
- name: E2E tests
  run: |
    npx playwright install --with-deps
    npm run test:e2e

G) Test-Hilfen
1. Testcontainers: spinnt MySQL/MariaDB on-the-fly, übergibt DSN an Tests.
2. Seed-Daten als SQL-Fixtures.
3. Für Scheduler-Tests: Clock-Interface injizieren, Fake-Clock verwenden.
4. Große Tabellen: Generatoren für Massendaten, um Streaming/Batches realistisch zu prüfen.
