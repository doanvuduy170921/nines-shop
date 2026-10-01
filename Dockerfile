# ---- build stage ----
FROM golang:1.25-alpine AS build
WORKDIR /src

# Copy module files first so dependency download is cached between builds.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/api

# ---- runtime stage ----
FROM alpine:3.22
# ca-certificates: needed for TLS calls (e.g. Gmail SMTP). tzdata: correct time zones.
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 app

WORKDIR /app

# 1. Copy file binary server từ build stage sang
COPY --from=build /out/server .

# 2. Copy thư mục migrations từ build stage sang đúng đường dẫn app đang gọi
COPY --from=build /src/internal/db/migrations ./internal/db/migrations

# Upload folder owned by the non-root user
RUN mkdir -p /app/uploads && chown -R app:app /app
RUN chown -R app:app /app/internal/db/migrations

USER app

EXPOSE 8080
ENTRYPOINT ["./server"]