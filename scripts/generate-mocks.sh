#!/bin/bash

# generate-mocks.sh - Script to generate mocks using Mockery v3
# This script generates mocks for all interfaces in the project

set -e

echo "Generating mocks using Mockery v3..."

# Ensure mockery is installed
if ! command -v mockery &> /dev/null; then
    echo "Mockery not found. Installing..."
    go install github.com/vektra/mockery/v2@latest
fi

# Clean up existing mocks
echo "Cleaning up existing mocks..."
rm -rf mocks/

# Generate mocks for repository interfaces
echo "Generating repository interface mocks..."
mockery --dir=internal/domain/repositories --all --output=mocks/repositories --case=underscore

# Generate mocks for service interfaces (if any in the future)
# echo "Generating service interface mocks..."
# mockery --dir=internal/domain/services --all --output=mocks/services --case=underscore

# Generate mocks for application service interfaces (if any in the future)
# echo "Generating application service interface mocks..."
# mockery --dir=cmd/api-server/internal/application --all --output=mocks/application --case=underscore

echo "Mock generation completed successfully!"

# Display generated mocks
echo "Generated mocks:"
find mocks/ -name "*.go" -type f | sort

# Check if mocks compile
echo "Verifying that generated mocks compile..."
go build ./mocks/...

echo "All mocks generated and verified successfully!"