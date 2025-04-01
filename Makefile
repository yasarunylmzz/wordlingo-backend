# .env dosyasındaki çevresel değişkenleri yükle
-include .env
export

# Veritabanı bilgilerini .env dosyasından al
DB_NAME := $(DB_NAME)
DB_USER := $(DB_USER)
DB_PASSWORD := $(DB_PASSWORD)
DB_HOST := $(DB_HOST)
DB_PORT := $(DB_PORT)

# Ortam değişkeni olarak PostgreSQL şifresini ayarla
export PGPASSWORD := $(DB_PASSWORD)

.PHONY: all clean createdb schema run_sql help

# Tüm adımları sırayla çalıştırmak için "all" hedefi
all: clean createdb schema run_sql

# Veritabanını silme işlemi
clean: ## Veritabanını siler
	@echo "PostgreSQL veritabanı siliniyor: $(DB_NAME)"
	psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "\
		SELECT pg_terminate_backend(pg_stat_activity.pid) \
		FROM pg_stat_activity \
		WHERE pg_stat_activity.datname = '$(DB_NAME)' \
		AND pid <> pg_backend_pid();"
	psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "DROP DATABASE IF EXISTS $(DB_NAME);"

# Veritabanını oluşturma işlemi
createdb: ## Veritabanını yeniden oluşturur
	@echo "PostgreSQL veritabanı oluşturuluyor: $(DB_NAME)"
	psql -U $(DB_USER) -h db -p $(DB_PORT) -tc "SELECT 1 FROM pg_database WHERE datname = '$(DB_NAME)'" | grep -q 1 || \
	psql -U $(DB_USER) -h db -p $(DB_PORT) -c "CREATE DATABASE $(DB_NAME);"

# Şema dosyasını yükleme işlemi
schema: ## Şema dosyasını uygular
	@echo "Şema dosyası uygulanıyor: db/schema/schema.sql"
	psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -d $(DB_NAME) -f db/schema/schema.sql

# /schema dizinindeki SQL dosyalarını sırasıyla çalıştırma işlemi
run_sql: ## /schema dizinindeki tüm SQL dosyalarını sırasıyla çalıştırır
	@echo "/schema dizinindeki tüm SQL dosyaları sıralı olarak çalıştırılıyor..."
	for file in db/schema/*.sql; do \
		if [ -f $$file ]; then \
			echo "Çalıştırılıyor: $$file"; \
			psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -d $(DB_NAME) -f $$file; \
		else \
			echo "Dosya bulunamadı: $$file"; \
		fi \
	done

# Mevcut hedefleri ve açıklamalarını gösteren help komutu
help: ## Bu yardım mesajını gösterir
	@echo "Kullanılabilir komutlar:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?##/ {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
