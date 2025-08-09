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

echo "Generating sqlc for ${SERVICE_NAME} service..."

# Define Go tool path and add it to PATH
GOPATH_BIN="$(go env GOPATH)/bin"
export PATH=$PATH:$GOPATH_BIN

echo "Checking and installing sqlc if not found..."

if ! command -v sqlc &> /dev/null; then
    echo "sqlc not found. Installing..."
    
    # Use go install to get the latest version of sqlc
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    
    echo "sqlc installed."
else
    echo "sqlc is already installed."
fi

# Check if the service directory exists
if [ ! -d "./apps/${SERVICE_NAME}" ]; then
    echo "Error: Service directory ./apps/${SERVICE_NAME} not found."
    exit 1
fi

# Navigate to the correct directory and run sqlc generate
# This assumes the sqlc.yaml configuration file is located in the service's postgres directory
cd "./apps/${SERVICE_NAME}/internal/infra/persistence/postgres"

sqlc generate

echo "sqlc: generated for ${SERVICE_NAME}"