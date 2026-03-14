set shell := ["bash", "-cu"]

go := env_var_or_default("GO", "go")
gocache := env_var_or_default("GOCACHE", "/tmp/go-build-cache")

default:
  @just --list

wire:
  @echo "Generating Wire code..."
  ~/go/bin/wire

mocks:
  @echo "Generating mocks in organized directories..."
  ~/go/bin/mockery

sqlc:
  @echo "Generating SQLC code..."
  ~/go/bin/sqlc generate

build: sqlc
  @echo "Building application..."
  GOCACHE={{gocache}} {{go}} build -o go-vbs ./cmd/go-vbs

run: build
  @echo "Starting application in development mode..."
  GO_ENV=development ./go-vbs

test: sqlc mocks
  @echo "Running tests..."
  GOCACHE={{gocache}} {{go}} test -v ./...

test-coverage: sqlc mocks
  @echo "Running tests with coverage..."
  GOCACHE={{gocache}} {{go}} test -v -cover ./...

clean:
  @echo "Cleaning build artifacts..."
  rm -f go-vbs go-vbs-wire

install-wire:
  @echo "Installing Wire..."
  {{go}} install github.com/google/wire/cmd/wire@latest

install-mockery:
  @echo "Installing Mockery..."
  {{go}} install github.com/vektra/mockery/v2@latest

install-sqlc:
  @echo "Installing SQLC..."
  {{go}} install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

install-goose:
  @echo "Installing Goose..."
  {{go}} install github.com/pressly/goose/v3/cmd/goose@latest

install-tools: install-wire install-mockery install-sqlc install-goose
  @echo "All development tools installed"

migrate-create name:
  @echo "Creating new migration..."
  ~/go/bin/goose -dir internal/repository/migrations create {{name}} sql

migrate-up:
  @echo "Applying migrations..."
  ~/go/bin/goose -dir internal/repository/migrations sqlite3 vbs.db up

migrate-down:
  @echo "Rolling back last migration..."
  ~/go/bin/goose -dir internal/repository/migrations sqlite3 vbs.db down

migrate-status:
  @echo "Migration status..."
  ~/go/bin/goose -dir internal/repository/migrations sqlite3 vbs.db status

swagger-ui:
  curl -L https://github.com/swagger-api/swagger-ui/archive/refs/tags/v5.24.1.tar.gz | tar -xzf -
  mkdir -p swagger-ui
  mv swagger-ui-5.24.1/dist/* swagger-ui/
  sed -i 's|https://petstore.swagger.io/v2/swagger.json|/docs/openapi.yaml|' swagger-ui/swagger-initializer.js
  rm -rf swagger-ui-5.24.1
