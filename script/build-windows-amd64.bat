@echo off
echo Building Freedom-Core for Windows 64-bit (x64)...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-H=windowsgui" -o freedom-core-windows-x64.exe ../cmd/server
echo Build complete! Output: freedom-core-windows-x64.exe
pause