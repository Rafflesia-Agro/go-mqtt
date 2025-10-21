#!/bin/bash

# Create build directory if it doesn't exist
mkdir -p build

# Build the Go application for the current platform
echo "Building server executable..."
go build -o build/server -ldflags="-s -w" .

echo "Server executable built successfully in build/server"