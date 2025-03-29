#!/bin/bash
rm -rf bin
# Windows 64 位
GOOS=windows GOARCH=amd64 go build -o bin/jilaidian_go.exe ../cmd/jilaidian/


# Linux 64 位
GOOS=linux GOARCH=amd64 go build -o bin/jilaidian_go ../cmd/jilaidian/