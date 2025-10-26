#!/bin/bash
echo "Building Freedom-Core for Linux 32-bit (x86)..."
export GOOS=linux
export GOARCH=386
go build -o freedom-core-linux-x86 ../cmd/server
echo "Build complete! Output: freedom-core-linux-x86"