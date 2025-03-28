#!/bin/bash
rm -rf bin
# Windows 64 位
GOOS=windows GOARCH=amd64 go build -o bin/jilaidian_go.exe


# Linux 64 位
GOOS=linux GOARCH=amd64 go build -o bin/jilaidian_go