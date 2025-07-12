#!/bin/bash
appName="yianqi_go"

echo $appName
rm -rf bin
go generate ../...

# Windows 64 位
if ! GOOS=windows GOARCH=amd64 go build -o bin/$appName.exe ../cmd/mainExe/;then
    echo "编译失败"
    exit 1
fi


# Linux 64 位
if ! GOOS=linux GOARCH=amd64 go build -o bin/$appName ../cmd/mainExe/; then
    echo "编译失败"
    exit 1
fi


if  [ "$1" = "up" ]; then
    echo "上传到七牛"
    python ./7niu/upload_7niu.py ./bin/$appName --filename=$appName
else
    echo "不进行上传"
fi