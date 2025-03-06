#!/bin/bash

# Build the Windows version of goLAPS
GOOS=windows GOARCH=amd64 go build -o bin/golaps.exe src/*.go