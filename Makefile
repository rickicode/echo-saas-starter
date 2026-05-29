.PHONY: build run test dev lint clean

# Build backend and frontend
build:
	cd ui && npm run build
	go build -o bin/server ./cmd/server

# Run the server
run:
	go run ./cmd/server

# Run all tests
test:
	go test ./...
	cd ui && npm run test

# Run development servers concurrently
dev:
	@echo "Starting backend and frontend dev servers..."
	@trap 'kill 0' EXIT; \
		(cd ui && npm run dev) & \
		go run ./cmd/server & \
		wait

# Run linters
lint:
	go vet ./...
	cd ui && npm run lint

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf ui/dist
