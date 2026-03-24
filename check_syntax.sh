#!/bin/bash

echo "🔍 Checking Go syntax..."

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Run go vet to check for common mistakes
echo "Running go vet..."
go vet ./...
if [ $? -ne 0 ]; then
    echo "❌ go vet found issues"
    exit 1
fi

# Try to build
echo "Building project..."
go build -o server ./cmd/server
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ All checks passed!"
echo "Binary created: ./server"
