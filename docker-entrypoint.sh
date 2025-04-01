#!/bin/bash
set -e

# Load .env file
source .env

# Wait for the database to be ready
./wait-for-it.sh db:5432 --timeout=30 --strict

# Check if the database exists
if PGPASSWORD=$DB_PASSWORD psql -h db -U postgres -d postgres -c '\l' | grep -q "flashcards"; then
  echo "Database exists, dropping it..."
  PGPASSWORD=$DB_PASSWORD psql -h db -U postgres -d postgres -c "DROP DATABASE IF EXISTS flashcards;" || { echo "Failed to drop the database"; exit 1; }
fi

# Create a new database
echo "Creating the new database..."
PGPASSWORD=$DB_PASSWORD make createdb || { echo "Failed to create the database"; exit 1; }

# Apply schema
echo "Applying schema..."
PGPASSWORD=$DB_PASSWORD make schema || { echo "Failed to apply schema"; exit 1; }

# Start the application
exec "./server"
