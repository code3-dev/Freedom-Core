#!/bin/bash
echo "Building Freedom-Core for Windows 32-bit (x86)..."
export GOOS=windows
export GOARCH=386
go build -ldflags "-H=windowsgui" -o freedom-core-windows-x86.exe ../cmd/server
echo "Build complete! Output: freedom-core-windows-x86.exe"