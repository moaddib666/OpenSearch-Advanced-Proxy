#!/bin/bash
set -e

# Function to check if ClickHouse server is ready
function clickhouse_ready(){
    clickhouse-client -q "SELECT 1" >/dev/null 2>&1
}

# Wait for ClickHouse Server to be ready
until clickhouse_ready; do
    echo "Waiting for ClickHouse server..."
    sleep 1
done

# Create database and tables
echo "Creating database and tables in ClickHouse..."

clickhouse-client -q "CREATE DATABASE IF NOT EXISTS $CLICKHOUSE_DB"

# Create table schema based on provided configuration
clickhouse-client -q "
CREATE TABLE IF NOT EXISTS $CLICKHOUSE_DB.example (
    datetime DateTime,
    message String
) ENGINE = MergeTree()
ORDER BY datetime"

clickhouse-client --query="INSERT INTO $CLICKHOUSE_DB.example FORMAT JSONEachRow" --input_format_skip_unknown_fields=1 < $INITIAL_DATA_PATH

echo "Database and tables created successfully."
