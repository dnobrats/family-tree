# ---------- Build stage ----------
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cài cert (để HTTPS / DB TLS nếu cần)
RUN apk add --no-cache ca-certificates

# Copy module files trước để cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o genealogy-be cmd/server/main.go

# ---------- Runtime stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/genealogy-be /app/genealogy-be

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/genealogy-be"]

