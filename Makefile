.PHONY: wire build run clean test mocks install-tools sqlc install-sqlc migrate migrate-create migrate-up migrate-down migrate-status install-goose

# Generate Wire dependency injection code
wire:
	@echo "Generating Wire code..."
	~/go/bin/wire

# Generate mocks using mockery (organized in mocks/ directories)
mocks:
	@echo "Generating mocks in organized directories..."
	~/go/bin/mockery

# Build the application
build:
	@echo "Building application..."
	go build -o go-vbs .

# Run the application in development mode
run: build
	@echo "Starting application in development mode..."
	GO_ENV=development ./go-vbs

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -cover ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f go-vbs go-vbs-wire

# Install Wire tool if not present
install-wire:
	@echo "Installing Wire..."
	go install github.com/google/wire/cmd/wire@latest

# Install mockery tool if not present
install-mockery:
	@echo "Installing Mockery..."
	go install github.com/vektra/mockery/v2@latest

# Install SQLC tool if not present
install-sqlc:
	@echo "Installing SQLC..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate SQLC code
sqlc:
	@echo "Generating SQLC code..."
	~/go/bin/sqlc generate

# Install Goose migration tool if not present
install-goose:
	@echo "Installing Goose..."
	go install github.com/pressly/goose/v3/cmd/goose@latest

# Create a new migration file
migrate-create:
	@echo "Creating new migration..."
	@read -p "Enter migration name: " name; \
	~/go/bin/goose -dir repository/migrations create $$name sql

# Apply all migrations
migrate-up:
	@echo "Applying migrations..."
	~/go/bin/goose -dir repository/migrations sqlite3 vbs.db up

# Rollback last migration
migrate-down:
	@echo "Rolling back last migration..."
	~/go/bin/goose -dir repository/migrations sqlite3 vbs.db down

# Show migration status
migrate-status:
	@echo "Migration status..."
	~/go/bin/goose -dir repository/migrations sqlite3 vbs.db status

# Install all development tools
install-tools: install-wire install-mockery install-sqlc install-goose
	@echo "All development tools installed"

# Download and setup Swagger UI
swagger-ui:
	curl -L https://github.com/swagger-api/swagger-ui/archive/refs/tags/v5.24.1.tar.gz | \
	tar -xzf - && \
	mkdir -p swagger-ui && \
	mv swagger-ui-5.24.1/dist/* swagger-ui/ && \
	rm -rf swagger-ui-5.24.1

# Help
help:
	@echo "Available commands:"
	@echo "  wire            - Generate Wire dependency injection code"
	@echo "  mocks           - Generate mocks using mockery"
	@echo "  sqlc            - Generate SQLC database code"
	@echo "  build           - Build the application (includes wire and sqlc generation)"
	@echo "  run             - Run the application in development mode"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  clean           - Clean build artifacts"
	@echo "  install-wire    - Install Wire tool"
	@echo "  install-mockery - Install Mockery tool"
	@echo "  install-sqlc    - Install SQLC tool"
	@echo "  install-goose   - Install Goose migration tool"
	@echo "  install-tools   - Install all development tools"
	@echo "  migrate-create  - Create a new migration file"
	@echo "  migrate-up      - Apply all pending migrations"
	@echo "  migrate-down    - Rollback last migration"
	@echo "  migrate-status  - Show migration status"
	@echo "  swagger-ui      - Download and setup Swagger UI (if not already installed)."
	@echo "                  	Remove the swagger-ui directory manually to update the Swagger UI."
	@echo "                  	In swagger-ui/swagger-initializer.js, change the url /docs/openapi.yaml"
	@echo "  help            - Show this help" 