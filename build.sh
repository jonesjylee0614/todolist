#!/bin/bash

# 编译目标: CentOS 7.8 (linux/amd64)
echo "开始编译 Go 项目到 CentOS 7.8..."

# 设置编译环境
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

# 编译信息
VERSION=$(git describe --tags --always 2>/dev/null || echo "unknown")
BUILD_TIME=$(date +%Y-%m-%d_%H:%M:%S)
LDFLAGS="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# 编译
echo "编译目标: ${GOOS}/${GOARCH}"
echo "版本信息: ${VERSION}"
echo "构建时间: ${BUILD_TIME}"

go build -ldflags="${LDFLAGS}" -o todolist-server main.go

if [ $? -eq 0 ]; then
    echo "编译成功!"
    echo "输出文件: todolist-server"
    echo "文件大小: $(ls -lh todolist-server | awk '{print $5}')"
    
    # 验证
    echo "文件类型: $(file todolist-server)"
    echo ""
    echo "部署步骤:"
    echo "1. 传输到服务器: scp todolist-server user@server:/path/"
    echo "2. 在服务器上运行: chmod +x todolist-server && ./todolist-server"
else
    echo "编译失败!"
    exit 1
fi
