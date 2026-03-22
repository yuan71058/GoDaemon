@echo off
echo Building GoDaemon 64-bit DLL...
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
go build -buildmode=c-shared -o damo64.dll ./exports
echo Build complete: damo64.dll
pause
