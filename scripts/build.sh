#!/bin/bash
rm -rf bin
# Windows 64 位
GOOS=windows GOARCH=amd64 go build -o bin/yilingshequ_go.exe ../cmd/mainExe/


# Linux 64 位
GOOS=linux GOARCH=amd64 go build -o bin/yilingshequ_go ../cmd/mainExe/
GOOS=linux GOARCH=amd64 go build -o bin/downloader ../cmd/download7niu/