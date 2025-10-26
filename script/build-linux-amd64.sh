#!/bin/bash
echo "Building Freedom-Core for Linux 64-bit (x64)..."
export GOOS=linux
export GOARCH=amd64
go build -o freedom-core-linux-x64 ../cmd/server
echo "Build complete! Output: freedom-core-linux-x64"