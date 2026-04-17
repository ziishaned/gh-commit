.PHONY: build test clean install

# Build the extension
build:
	@echo "Building gh-commit..."
	@mkdir -p bin
	@go build -o bin/commit ./cmd/commit
	@echo "✅ Built bin/commit"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin
	@echo "✅ Cleaned"

# Install as gh extension
install: build
	@echo "Installing as gh extension..."
	@gh extension install .
	@echo "✅ Installed"

# Uninstall extension
uninstall:
	@echo "Uninstalling gh extension..."
	@gh extension remove commit
	@echo "✅ Uninstalled"
