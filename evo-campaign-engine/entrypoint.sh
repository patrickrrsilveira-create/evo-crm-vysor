#!/bin/sh

DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}&x-migrations-table=evo_campaign_engine_schema_migrations"

echo "🚀 Starting EVO Campaign Engine..."
echo "📊 Database: ${DB_HOST}:${DB_PORT}/${DB_NAME}"
echo "🔄 Running migrations..."

./migrate -database "$DB_URL" -path ./migrations up

if [ $? -eq 0 ]; then
    echo "✅ Migrations completed successfully!"
else
    echo "❌ Migration failed!"
    exit 1
fi

echo "🚀 Starting application..."
exec ./campaign-engine
