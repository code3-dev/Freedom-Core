@echo off
echo Building Freedom-Core for Windows 32-bit (x86)...
set GOOS=windows
set GOARCH=386
go build -ldflags "-H=windowsgui" -o freedom-core-windows-x86.exe ../cmd/server
echo Build complete! Output: freedom-core-windows-x86.exe
pause