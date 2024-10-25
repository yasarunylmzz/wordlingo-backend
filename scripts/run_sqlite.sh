#!/bin/bash

# JSON dosyasından veritabanı bağlantı bilgilerini al
DB_NAME=$(jq -r ".db_name" config.json)
DB_USER=$(jq -r ".db_user" config.json)
DB_PASSWORD=$(jq -r ".db_password" config.json)
DB_HOST=$(jq -r ".db_host" config.json)
DB_PORT=$(jq -r ".db_port" config.json)

# Şifreyi ortam değişkeni olarak ayarla
export PGPASSWORD="$DB_PASSWORD"

# Veritabanını sil
echo "PostgreSQL veritabanı siliniyor: $DB_NAME"
psql -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -c "DROP DATABASE IF EXISTS $DB_NAME;"

# Veritabanını yeniden oluştur
echo "PostgreSQL veritabanı yeniden oluşturuluyor: $DB_NAME"
psql -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -c "CREATE DATABASE $DB_NAME;"

# Şema dosyasını çalıştırarak veritabanı şemasını yükle
echo "Şema dosyası uygulanıyor: db/schema/schema.sql"
psql -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -f "db/schema/schema.sql"

# /schema dizinindeki types.sql ve tables.sql dosyalarını sıralı olarak çalıştır
echo "/schema dizinindeki types.sql ve tables.sql dosyaları sıralı olarak çalıştırılıyor..."
for file in db/schema/*.sql; do
    if [ -f "$file" ]; then
        echo "Çalıştırılıyor: $file"
        psql -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -f "$file"
    else
        echo "Dosya bulunamadı: $file"
    fi
done

# PGPASSWORD ortam değişkenini sıfırla (güvenlik için)
unset PGPASSWORD

echo "PostgreSQL veritabanı başarıyla yeniden başlatıldı ve tüm SQL dosyaları uygulandı."
