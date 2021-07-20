@echo off

where go >NUL

if %ERRORLEVEL% == 1 (
  echo "go is not installed"
  exit /b 1
)

go env -w GOOS=js GOARCH=wasm
pushd denops\go
go build -o main.wasm .
popd
go env -u GOOS GOARCH
