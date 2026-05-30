.PHONY: build run test dev lint clean docker-build docker-run docker-up docker-down frontend-build frontend-dev

# Build backend and frontend
build: frontend-build
	go build -o bin/server ./cmd/server

# Run the server
run:
	go run ./cmd/server

# Run all tests
test:
	go test ./... -v

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

# Clean build artifacts
clean:
	rm -rf bin/ ui/dist

# Docker commands
docker-build:
	docker build -t echo-saas-starter .

docker-run:
	docker run -p 8080:8080 echo-saas-starter

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Frontend commands
frontend-build:
	cd ui && npm run build

frontend-dev:
	cd ui && npm run dev
