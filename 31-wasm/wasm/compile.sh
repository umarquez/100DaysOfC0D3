#!/bin/bash

source /etc/profile
GOOS=js GOARCH=wasm go build -o ../www/app.wasm .
