#!/bin/bash


# Export the go binary
# export GOPATH=/Users/$USER/go

cp $(go env GOROOT)/misc/wasm/wasm_exec.js .
GOOS=js GOARCH=wasm go build -o main.wasm


# Then, run an http server, for example:
#   python3 -m http.server 5000
# or:
#   npm i -g serve
#   serve
