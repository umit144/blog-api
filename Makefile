# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application using gotestsum
test:
	@echo "Running tests..."
	@if command -v gotestsum > /dev/null; then \
		gotestsum --format pkgname -- ./tests/... -v; \
	else \
		read -p "gotestsum is not installed. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install gotest.tools/gotestsum@latest; \
			gotestsum --format pkgname -- ./tests/... -v; \
		else \
			echo "Running tests without gotestsum..."; \
			go test ./tests/... -v; \
		fi; \
	fi

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

migration:
	@goose -dir internal/database/migration postgres "host=localhost port=5432 dbname=postgres user=postgres password=postgres sslmode=disable" up

.PHONY: all build run test clean