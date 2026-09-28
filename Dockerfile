# ---- build stage ----
# Use the same Go version as in go.mod (e.g. golang:1.25-alpine).
FROM golang:1.25-alpine AS build
WORKDIR /src

# Copy module files first so dependency download is cached between builds.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Change "." if main.go is in a sub folder (e.g. ./cmd/server).
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/api

# ---- runtime stage ----
FROM alpine:3.22
# ca-certificates: needed for TLS calls (e.g. Gmail SMTP). tzdata: correct time zones.
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 app

WORKDIR /app
COPY --from=build /out/server .

# Upload folder owned by the non-root user. With a named volume, Docker copies
# this ownership on first mount; with a host folder you must chown it to 10001.
RUN mkdir -p /app/uploads && chown -R app:app /app
USER app

EXPOSE 8080
ENTRYPOINT ["./server"]