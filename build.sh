#!/bin/bash

# Memberi tahu Go untuk tidak menggunakan library C (CGO_ENABLED=0)
# dan secara spesifik membuat build untuk OS Linux (GOOS=linux)
echo "Building static Linux server executable..."
CGO_ENABLED=0 GOOS=linux go build -o server -ldflags="-s -w" .

echo "Server executable built successfully in ./server"