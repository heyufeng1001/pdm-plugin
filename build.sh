#!/bin/bash
RUN_NAME=hertz_service
mkdir -p output/bin
mkdir -p output/conf
cp script/* output 2>/dev/null
cp conf/* output/conf 2>/dev/null
chmod +x output/bootstrap.sh
chmod +x output/destroy.sh
# 自行适配运行环境和编译变量
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o output/bin/${RUN_NAME}