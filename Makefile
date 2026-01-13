.PHONY: build run test clean docker-build docker-up docker-down docker-logs

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/server ./cmd/server
	@echo "Build complete: bin/server"

# Run the application
run:
	@echo "Running application..."
	@go run ./cmd/server

# Run tests
test:
	@echo "Running tests..."
	@go test -v -cover ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t rate-limiter-service:latest .
	@echo "Docker image built: rate-limiter-service:latest"

# Run with Docker Compose (up)
docker-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d
	@echo "Services started. Access:"
	@echo "  - API: http://localhost:8080"
	@echo "  - Prometheus: http://localhost:9090"
	@echo "  - Grafana: http://localhost:3000 (admin/admin)"

# Run with Docker Compose in foreground
docker-up-fg:
	@docker-compose up

# Stop Docker Compose
docker-down:
	@echo "Stopping services..."
	@docker-compose down
	@echo "Services stopped"

# View Docker Compose logs
docker-logs:
	@docker-compose logs -f

# View logs for specific service
docker-logs-service:
	@docker-compose logs -f rate-limiter

# Restart services
docker-restart:
	@echo "Restarting services..."
	@docker-compose restart
	@echo "Services restarted"

