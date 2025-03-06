#!/bin/bash

# Build the Linux version of goLAPS
GOOS=linux GOARCH=amd64 go build -o bin/golaps-amd64 src/*.go