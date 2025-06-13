.PHONY: wire build run clean test mocks install-tools

# Generate Wire dependency injection code
wire:
	@echo "Generating Wire code..."
	~/go/bin/wire

# Generate mocks using mockery (organized in mocks/ directories)
mocks:
	@echo "Generating mocks in organized directories..."
	~/go/bin/mockery

# Build the application
build: wire
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

# Install all development tools
install-tools: install-wire install-mockery
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
	@echo "  wire          - Generate Wire dependency injection code"
	@echo "  mocks         - Generate mocks using mockery"
	@echo "  build         - Build the application (includes wire generation)"
	@echo "  run           - Run the application in development mode"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  install-wire  - Install Wire tool"
	@echo "  install-mockery - Install Mockery tool"
	@echo "  install-tools - Install all development tools"
	@echo "  swagger-ui    - Download and setup Swagger UI (if not already installed)."
	@echo "                  Remove the swagger-ui directory manually to update the Swagger UI."
	@echo "                  In swagger-ui/swagger-initializer.js, change the url /docs/openapi.yaml"
	@echo "  help          - Show this help" 