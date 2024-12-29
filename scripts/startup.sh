#!/bin/sh

# Function to log with timestamp
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1"
}

# Function to check migration status
check_migration_status() {
    log "Checking current migration version..."
    migrate -path db/migrations -database "$1" version
    if [ $? -eq 0 ]; then
        log "Successfully retrieved migration version"
    else
        log "WARNING: Could not retrieve migration version"
    fi
}

# Function to run database migrations
run_migrations() {
    log "Starting database migration process..."

    # Load environment variables from .env 
    if [ -f ".env" ]; then
        log "Loading environment variables from .env"
        export $(grep -v '^#' .env | xargs) 
    else
        log "WARNING: .env file not found. Using system environment variables."
    fi

    # Check required environment variables
    for var in DB_HOST DB_USER DB_PASSWORD DB_NAME; do
        if [ -z "$(eval echo \$$var)" ]; then
            log "ERROR: Required environment variable $var is not set"
            exit 1
        fi
    done

    log "Checking database connection..."
    max_retries=30
    counter=0

    while [ $counter -lt $max_retries ]; do
        if pg_isready -h "${DB_HOST}" -U "${DB_USER}"; then
            log "Database connection established successfully"
            break
        fi
        log "Attempt $((counter + 1))/$max_retries: Database not ready, waiting..."
        counter=$((counter + 1))
        sleep 2
    done

    if [ $counter -eq $max_retries ]; then
        log "ERROR: Database connection timeout after $max_retries attempts"
        exit 1
    fi

    # Construct database URL
    dbURL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}/${DB_NAME}?sslmode=disable"

    # Check current migration status
    log "Checking migration status before applying migrations..."
    check_migration_status "$dbURL"

    # Run migrations
    log "Applying database migrations..."
    migrate -path db/migrations -database "$dbURL" up

    if [ $? -eq 0 ]; then
        log "Migrations completed successfully"
        log "Final migration status:"
        check_migration_status "$dbURL"
    else
        log "ERROR: Migration failed"
        log "Current migration status:"
        check_migration_status "$dbURL"
        exit 1
    fi
}

# Main execution flow
log "=== Starting Identity Service Initialization ==="

log "Step 1/2: Database Migration"
run_migrations

log "Step 2/2: Starting Application"
log "=== Identity Service Initialization Complete ==="
exec ./main