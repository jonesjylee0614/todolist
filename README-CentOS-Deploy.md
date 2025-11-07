# Go 项目 CentOS 7.8 部署指南

本项目提供了在 Windows 11 WSL 环境中编译 Go 项目并部署到 CentOS 7.8 的完整解决方案。

## 快速开始

### 方法一：使用 Windows 批处理文件（推荐）

```cmd
# 在项目根目录的 backend 文件夹中运行
cd backend
build.bat
```

### 方法二：使用 WSL 脚本

```bash
# 在 WSL 环境中运行
cd /mnt/d/WorkSpace/code/2025/todolist/backend
chmod +x build.sh
./build.sh
```

### 方法三：手动编译

```cmd
# Windows CMD
cd backend
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags="-w -s" -o todolist-server main.go
```

## 文件说明

- `build.bat` - Windows 批处理编译脚本
- `build.sh` - Linux/WSL 编译脚本
- `deploy.sh` - 自动化部署脚本
- `todolist.service` - systemd 服务配置文件
- `build-centos.md` - 详细的技术文档

## 部署步骤

### 1. 编译项目

选择上述任一方法编译项目，生成 `todolist-server` 二进制文件。

### 2. 部署到服务器

```bash
# 使用部署脚本（需要 SSH 访问权限）
./deploy.sh root@your-server-ip /opt/todolist

# 或手动部署
scp todolist-server user@server:/opt/todolist/
ssh user@server
chmod +x /opt/todolist/todolist-server
```

### 3. 配置系统服务

```bash
# 复制服务文件
sudo cp todolist.service /etc/systemd/system/

# 重载 systemd
sudo systemctl daemon-reload

# 启用服务
sudo systemctl enable todolist

# 启动服务
sudo systemctl start todolist

# 查看状态
sudo systemctl status todolist
```

## 环境要求

### 开发环境（Windows 11 + WSL）
- Go 1.24.0+
- Git
- WSL2（推荐 Ubuntu 20.04/22.04）

### 目标环境（CentOS 7.8）
- glibc 2.17+
- systemd（用于服务管理）
- 64位 x86_64 架构

## 编译参数说明

- `GOOS=linux` - 目标操作系统为 Linux
- `GOARCH=amd64` - 目标架构为 64位 x86
- `CGO_ENABLED=0` - 禁用 CGO，生成静态链接的二进制文件
- `-ldflags="-w -s"` - 去除调试信息和符号表，减小文件大小

## 常见问题

### 1. 编译失败
- 确保 Go 版本 >= 1.24.0
- 检查网络连接，确保能下载依赖
- 清理缓存：`go clean -modcache`

### 2. 运行时错误
- 检查文件权限：`chmod +x todolist-server`
- 查看 SELinux 状态：`sestatus`
- 检查端口占用：`netstat -tlnp | grep :8080`

### 3. 服务启动失败
- 查看日志：`sudo journalctl -u todolist -f`
- 检查配置文件路径
- 确认数据库连接配置

## 性能优化

### 减小二进制文件大小

```bash
# 使用 upx 压缩（可选）
upx --best todolist-server
```

### 生产环境配置

```bash
# 设置环境变量
export GIN_MODE=release
export PORT=8080
```

## 监控和维护

### 查看服务状态
```bash
sudo systemctl status todolist
sudo journalctl -u todolist -f
```

### 重启服务
```bash
sudo systemctl restart todolist
```

### 更新部署
```bash
# 停止服务
sudo systemctl stop todolist

# 替换二进制文件
cp new-todolist-server /opt/todolist/bin/todolist-server

# 启动服务
sudo systemctl start todolist
```

## 安全建议

1. 使用非 root 用户运行服务
2. 配置防火墙规则
3. 使用 HTTPS（配置反向代理）
4. 定期更新系统和依赖
5. 配置日志轮转

## 技术支持

如遇到问题，请检查：
1. 编译环境配置
2. 目标系统兼容性
3. 网络连接
4. 权限设置

更多详细信息请参考 `build-centos.md` 文档。
