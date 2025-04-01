# Build aşaması
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Go modlarını indir
COPY go.mod go.sum ./
RUN go mod download

# Uygulama kodunu kopyala
COPY . .

# Uygulamayı derle
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./server.go


# Runtime aşaması
FROM alpine:3.18.0

# Gerekli paketleri yükle
RUN apk --no-cache add make bash postgresql-client ca-certificates jq

WORKDIR /app

# Builder'dan sadece gerekli dosyaları kopyala
COPY --from=builder /app/server .
# COPY --from=builder /app/.env .
COPY --from=builder /app/wait-for-it.sh .
COPY --from=builder /app/Makefile .
COPY --from=builder /app/db ./db

# Gerekli izinleri ayarla
RUN chmod +x wait-for-it.sh

# Özel entrypoint script'ini kopyala ve izinlerini ayarla
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]
