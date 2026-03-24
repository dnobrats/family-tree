.PHONY: build run dev test clean

# Build binary
build:
	go build -o server ./cmd/server

# Run server
run: build
	./server

# Development mode với hot reload (cần air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f server

# Format code
fmt:
	go fmt ./...

# Run linter (cần golangci-lint)
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

# Generate password hash (sử dụng genpass.go)
genpass:
	go run genpass.go

# Reset database sequences (sau khi xóa records)
reset-sequences:
	@echo "🔄 Resetting database sequences..."
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f migrations/reset_all_sequences.sql
	@echo "✅ Done!"
