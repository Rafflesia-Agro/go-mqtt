#!/bin/bash

# Build the Go application for the current platform
echo "Building server executable..."
go build -o server -ldflags="-s -w" .

echo "Server executable built successfully in ./server"