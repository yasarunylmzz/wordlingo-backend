#!/bin/bash
set -e

# Veritabanı hazır olana kadar bekle
./wait-for-it.sh db:5432 --timeout=30 --strict

# Veritabanı var mı kontrol et
if ! PGPASSWORD=abc123 psql -h db -U postgres -d postgres -c '\l' | grep -q "flashcards"; then
  echo "Veritabanı bulunamadı, oluşturuluyor..."
  PGPASSWORD=abc123 make createdb || { echo "Veritabanı oluşturulamadı"; exit 1; }
else
  echo "Veritabanı mevcut."
fi

# Şemalar var mı kontrol et
if ! PGPASSWORD=abc123 psql -h db -U postgres -d flashcards -c '\dt' | grep -q "users"; then
  echo "Şemalar yok, uygulama yapılıyor..."
  PGPASSWORD=abc123 make schema || { echo "Şemalar uygulanamadı"; exit 1; }
else
  echo "Şemalar zaten mevcut."
fi

# Uygulamayı başlat
exec "./server"
