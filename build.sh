#!/usr/bin/env bash
mkdir -p output/

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o output/main 