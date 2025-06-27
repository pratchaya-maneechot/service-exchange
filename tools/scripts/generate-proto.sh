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

# Define Go tool paths within the script
GOPATH_BIN="$(go env GOPATH)/bin"
PROTOC_GEN_GO="${GOPATH_BIN}/protoc-gen-go"
PROTOC_GEN_GRPC="${GOPATH_BIN}/protoc-gen-go-grpc"

echo "Generating proto for ${SERVICE_NAME} service..."

# Use the full paths
protoc --plugin=protoc-gen-go="${PROTOC_GEN_GO}" \
    --go_out=. --go_opt=paths=source_relative \
    --plugin=protoc-gen-go-grpc="${PROTOC_GEN_GRPC}" \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./apps/${SERVICE_NAME}/api/proto/**/*.proto

echo "proto: generated for ${SERVICE_NAME}"