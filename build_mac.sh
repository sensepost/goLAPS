#!/bin/bash

# Build the Linux version of goLAPS
GOOS=darwin GOARCH=arm64 go build -o bin/golaps-arm64 src/*.go