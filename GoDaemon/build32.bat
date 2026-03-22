@echo off
echo Building GoDaemon 32-bit DLL...
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=386
go build -buildmode=c-shared -o godaemon32.dll ./exports
echo Build complete: godaemon32.dll
pause
