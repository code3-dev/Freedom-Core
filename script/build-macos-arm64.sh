#!/bin/bash
echo "Building Freedom-Core for macOS Apple Silicon (ARM64)..."
export GOOS=darwin
export GOARCH=arm64
go build -o freedom-core-macos-arm64 ../cmd/server
echo "Build complete! Output: freedom-core-macos-arm64"