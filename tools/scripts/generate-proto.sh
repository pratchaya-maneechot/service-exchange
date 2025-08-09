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

echo "Checking and installing Go proto plugins..."

go_plugins_to_install=(
    "google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
)

for plugin in "${go_plugins_to_install[@]}"; do
    plugin_name=$(basename "$plugin" | cut -d'@' -f1)
    if ! command -v "$plugin_name" &> /dev/null; then
        echo "Installing $plugin_name..."
        go install "$plugin"
    else
        echo "$plugin_name is already installed."
    fi
done

echo "Go proto plugins are ready."

# Define Go tool paths within the script
GOPATH_BIN="$(go env GOPATH)/bin"
PROTOC_GEN_GO="${GOPATH_BIN}/protoc-gen-go"
PROTOC_GEN_GRPC="${GOPATH_BIN}/protoc-gen-go-grpc"

echo "Checking and installing protoc if not found..."

# Check if protoc is already installed
if ! command -v protoc &> /dev/null; then
    echo "protoc not found. Downloading and installing..."

    # Determine OS and architecture
    OS="$(uname -s)"
    ARCH="$(uname -m)"

    PROTOC_VERSION="27.1"
    DOWNLOAD_URL=""
    INSTALL_DIR="/usr/local"

    if [[ "$OS" == "Linux" ]]; then
        if [[ "$ARCH" == "x86_64" ]]; then
            DOWNLOAD_URL="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip"
        elif [[ "$ARCH" == "aarch64" ]]; then
            DOWNLOAD_URL="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-aarch_64.zip"
        fi
    elif [[ "$OS" == "Darwin" ]]; then # macOS
        if [[ "$ARCH" == "x86_64" ]]; then
            DOWNLOAD_URL="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-osx-x86_64.zip"
        elif [[ "$ARCH" == "arm64" ]]; then
            DOWNLOAD_URL="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-osx-aarch_64.zip"
        fi
    fi

    if [[ -z "$DOWNLOAD_URL" ]]; then
        echo "Error: Unsupported OS or architecture."
        exit 1
    fi
    
    # Download and unzip protoc
    curl -L -o /tmp/protoc.zip "${DOWNLOAD_URL}"
    sudo unzip /tmp/protoc.zip -d "${INSTALL_DIR}"
    rm /tmp/protoc.zip
    
    echo "protoc installed to ${INSTALL_DIR}"
else
    echo "protoc is already installed."
fi

echo "Generating proto for ${SERVICE_NAME} service..."

# Use the full paths
protoc --plugin=protoc-gen-go="${PROTOC_GEN_GO}" \
    --go_out=. --go_opt=paths=source_relative \
    --plugin=protoc-gen-go-grpc="${PROTOC_GEN_GRPC}" \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./apps/${SERVICE_NAME}/api/proto/**/*.proto

echo "proto: generated for ${SERVICE_NAME}"