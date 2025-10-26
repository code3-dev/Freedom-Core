#!/bin/bash
echo "Building Freedom-Core for macOS Intel (x64)..."
export GOOS=darwin
export GOARCH=amd64
go build -o freedom-core-macos-x64 ../cmd/server
echo "Build complete! Output: freedom-core-macos-x64"