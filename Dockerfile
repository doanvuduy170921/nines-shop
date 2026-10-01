# ---- runtime stage ----
FROM alpine:3.22
# ca-certificates: needed for TLS calls (e.g. Gmail SMTP). tzdata: correct time zones.
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 app

WORKDIR /app

# 1. Copy file binary server từ build stage sang
COPY --from=build /out/server .

# 2. BỔ SUNG: Copy thư mục migrations từ build stage sang đúng cấu trúc đường dẫn mà code Go đang gọi
COPY --from=build /src/internal/db/migrations ./internal/db/migrations

# Upload folder owned by the non-root user. With a named volume, Docker copies
# this ownership on first mount; with a host folder you must chown it to 10001.
RUN mkdir -p /app/uploads && chown -R app:app /app

# Phân quyền cho user app có thể đọc được thư mục migrations (nếu cần)
RUN chown -R app:app /app/internal/db/migrations

USER app

EXPOSE 8080
ENTRYPOINT ["./server"]