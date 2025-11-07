# Windows 11 WSL 环境编译 Go 项目到 CentOS 7.8 指南

## 环境要求

### Windows 11 WSL 环境
- 已安装 WSL2 (推荐 Ubuntu 20.04 或 22.04)
- Go 1.24.0 或更高版本

### 目标环境 (CentOS 7.8)
- glibc 2.17 或更高版本
- 64位 x86_64 架构

## 编译步骤

### 1. 确认 WSL 中的 Go 版本

```bash
# 在 WSL 中执行
go version
# 应该显示 go1.24.0 或更高版本
```

### 2. 设置环境变量

```bash
# 设置目标操作系统和架构
export GOOS=linux
export GOARCH=amd64

# 设置 CGO 禁用（避免依赖系统库）
export CGO_ENABLED=0
```

### 3. 编译项目

```bash
# 进入项目目录
cd /mnt/d/WorkSpace/code/2025/todolist/backend

# 下载依赖
go mod download
go mod tidy

# 编译生产版本
go build -ldflags="-w -s" -o todolist-server main.go

# 查看编译结果
file todolist-server
# 应该显示: todolist-server: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked
```

### 4. 验证编译结果

```bash
# 检查二进制文件信息
ldd todolist-server
# 应该显示: not a dynamic executable

# 查看文件大小
ls -lh todolist-server
```

## 部署到 CentOS 7.8

### 1. 传输文件到服务器

```bash
# 使用 scp 传输
scp todolist-server user@your-centos-server:/path/to/deploy/

# 或使用 rsync
rsync -avz todolist-server user@your-centos-server:/path/to/deploy/
```

### 2. 在 CentOS 7.8 上运行

```bash
# 登录到 CentOS 服务器
ssh user@your-centos-server

# 给执行权限
chmod +x /path/to/deploy/todolist-server

# 运行服务
./todolist-server
```

## 编译脚本

创建一个自动化编译脚本 `build.sh`:

```bash
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
else
    echo "编译失败!"
    exit 1
fi
```

## 常见问题解决

### 1. CGO 相关问题

如果遇到 CGO 相关错误，确保设置：
```bash
export CGO_ENABLED=0
```

### 2. 依赖库问题

确保所有依赖都是纯 Go 实现，避免 C 依赖。如果必须使用 CGO，需要在目标系统上安装相应的开发库。

### 3. 运行时权限问题

在 CentOS 上运行时如果遇到权限问题：
```bash
# 检查 SELinux 状态
sestatus

# 临时禁用 SELinux（仅测试环境）
setenforce 0

# 或设置正确的 SELinux 上下文
chcon -t bin_t ./todolist-server
```

### 4. 端口占用

确保目标端口未被占用：
```bash
# 检查端口占用
netstat -tlnp | grep :8080
# 或
ss -tlnp | grep :8080
```

## 优化建议

### 1. 减小二进制文件大小

使用 upx 压缩（可选）：
```bash
# 安装 upx
sudo apt install upx-ucl

# 压缩二进制文件
upx --best todolist-server
```

### 2. 创建 systemd 服务

在 CentOS 7.8 上创建系统服务：

```ini
# /etc/systemd/system/todolist.service
[Unit]
Description=TodoList Server
After=network.target

[Service]
Type=simple
User=todolist
WorkingDirectory=/opt/todolist
ExecStart=/opt/todolist/todolist-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl enable todolist
sudo systemctl start todolist
sudo systemctl status todolist
```

## 验证部署

```bash
# 检查服务状态
curl -I http://localhost:8080/health

# 查看日志
sudo journalctl -u todolist -f
