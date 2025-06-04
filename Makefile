.PHONY: wire build run clean test

# Generate Wire dependency injection code
wire:
	@echo "Generating Wire code..."
	~/go/bin/wire

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

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f go-vbs go-vbs-wire

# Install Wire tool if not present
install-wire:
	@echo "Installing Wire..."
	go install github.com/google/wire/cmd/wire@latest

# Help
help:
	@echo "Available commands:"
	@echo "  wire         - Generate Wire dependency injection code"
	@echo "  build        - Build the application (includes wire generation)"
	@echo "  run          - Run the application in development mode"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  install-wire - Install Wire tool"
	@echo "  help         - Show this help" 