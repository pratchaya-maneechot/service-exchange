#!/bin/bash

# Ensure script exits if any command fails
set -e

# Arguments:
# 1: SERVICE_NAME - The name of the service (e.g., users, tasks)

SERVICE_NAME=$1

if [ -z "$SERVICE_NAME" ]; then
    echo "Error: Service name not provided."
    echo "Usage: $0 <service_name>"
    exit 1
fi

echo "Running migrations for ${SERVICE_NAME} service..."

GOPATH_BIN="$(go env GOPATH)/bin"
export PATH=$PATH:$GOPATH_BIN

echo "Checking and installing 'migrate' CLI if not found..."
if ! command -v migrate &> /dev/null; then
    echo "'migrate' CLI not found. Installing..."
    
    # Use go install to get the latest version of migrate with postgres driver
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    
    echo "'migrate' CLI installed."
else
    echo "'migrate' CLI is already installed."
fi

# Check if DB_URL is set
if [ -z "$DB_URL" ]; then
    echo "Error: DB_URL environment variable is not set."
    exit 1
fi

# Run migrate up command
migrate -path "./apps/${SERVICE_NAME}/internal/infra/persistence/postgres/migrations" \
        -database "${DB_URL}/${SERVICE_NAME}?sslmode=disable" \
        up

echo "Migrations for ${SERVICE_NAME} completed."