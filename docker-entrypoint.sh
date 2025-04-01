#!/bin/bash
set -e

# Ortam değişkenini kontrol et (development ya da production)
if [ -z "$GO_ENV" ]; then
  echo "GO_ENV environment variable is not set. Defaulting to 'production'."
  GO_ENV="production"
fi

# Ortama göre .env dosyasını yükle
echo "Loading .env.${GO_ENV} file"
source .env.${GO_ENV}

# Veritabanının hazır olmasını bekle
./wait-for-it.sh db:5432 --timeout=30 --strict

# Veritabanının mevcut olup olmadığını kontrol et
if PGPASSWORD=$DB_PASSWORD psql -h db -U postgres -d postgres -c '\l' | grep -q "$DB_NAME"; then
  echo "Database exists, dropping it..."
  PGPASSWORD=$DB_PASSWORD psql -h db -U postgres -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;" || { echo "Failed to drop the database"; exit 1; }
fi

# Yeni veritabanını oluştur
echo "Creating the new database..."
PGPASSWORD=$DB_PASSWORD make createdb || { echo "Failed to create the database"; exit 1; }

# Şemayı uygula
echo "Applying schema..."
PGPASSWORD=$DB_PASSWORD make schema || { echo "Failed to apply schema"; exit 1; }

# Uygulamayı başlat
echo "Starting the application..."
exec "./server"
