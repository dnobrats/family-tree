# syntax=docker/dockerfile:1

# ---------- build ----------
FROM golang:1.24-alpine AS build
WORKDIR /src

# Copy the whole module (small) and resolve dependencies with a warm cache.
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/createadmin ./cmd/createadmin

# ---------- runtime ----------
# distroless/static is tiny and has no shell — ideal for a static Go binary.
# Templates and static assets are embedded in the binary.
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/server /app/server
COPY --from=build /out/createadmin /app/createadmin
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
