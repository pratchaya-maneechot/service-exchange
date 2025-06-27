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

# Define Go tool path within the script
GOPATH_BIN="$(go env GOPATH)/bin"
WIRE_CMD="${GOPATH_BIN}/wire"

echo "Generating wire for ${SERVICE_NAME} service..."

# Change directory to the service's root before running wire
cd ./apps/${SERVICE_NAME}

# Use the full path
"${WIRE_CMD}" ./internal

echo "wire: generated for ${SERVICE_NAME}"