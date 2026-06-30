#!/bin/sh

# Construct the database URL with custom schema migrations table
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}&x-migrations-table=evo_core_community_schema_migrations"

echo "🚀 Starting EvoAI Core Service..."
echo "📊 Database: ${DB_HOST}:${DB_PORT}/${DB_NAME}"
echo "🔄 Running migrations..."

# Run migrations
./migrate -database "$DB_URL" -path ./migrations up

if [ $? -eq 0 ]; then
    echo "✅ Migrations completed successfully!"
else
    echo "❌ Migration failed!"
    exit 1
fi

echo "🚀 Starting application..."
exec ./main
