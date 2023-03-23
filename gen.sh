#!/usr/bin/zsh

docker run --rm --name protoc --mount type=bind,source=$(pwd)/proto,target=/app/proto registry.cn-qingdao.aliyuncs.com/ppg007/protoc-gen:latest /sbin/my_init -- bash gen.sh
