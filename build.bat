chcp 65001

@echo off
echo Building PortFinder for multiple platforms...

echo Building Windows version...
set CGO_ENABLED=0
go build -ldflags "-s -w" -o pf.exe

echo Building Linux version...
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags "-s -w -extldflags '-static'" -o pf-linux

echo Building macOS version...
set GOOS=darwin
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags "-s -w" -o pf-macos

echo Building macOS ARM64 version...
set GOOS=darwin
set GOARCH=arm64
set CGO_ENABLED=0
go build -ldflags "-s -w" -o pf-macos-arm64

echo Build completed!
echo Generated files:
echo   pf.exe (Windows)
echo   pf-linux (Linux)
echo   pf-macos (macOS Intel)
echo   pf-macos-arm64 (macOS Apple Silicon) 