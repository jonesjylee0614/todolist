#!/bin/bash

# 部署脚本 - 将编译好的二进制文件部署到 CentOS 7.8
# 使用方法: ./deploy.sh <user@server> <deploy_path>

if [ $# -lt 2 ]; then
    echo "使用方法: $0 <user@server> <deploy_path>"
    echo "示例: $0 root@192.168.1.100 /opt/todolist"
    exit 1
fi

SERVER=$1
DEPLOY_PATH=$2
BINARY="todolist-server"
SERVICE_FILE="todolist.service"

echo "开始部署到 ${SERVER}..."

# 检查二进制文件是否存在
if [ ! -f "$BINARY" ]; then
    echo "错误: 找不到编译好的二进制文件 $BINARY"
    echo "请先运行 ./build.sh 编译项目"
    exit 1
fi

# 创建部署目录
echo "创建部署目录..."
ssh $SERVER "sudo mkdir -p $DEPLOY_PATH/{bin,logs,config}"

# 传输二进制文件
echo "传输二进制文件..."
scp $BINARY $SERVER:$DEPLOY_PATH/bin/

# 传输服务文件
if [ -f "$SERVICE_FILE" ]; then
    echo "传输 systemd 服务文件..."
    scp $SERVICE_FILE $SERVER:/tmp/
fi

# 在服务器上执行设置
echo "在服务器上设置..."
ssh $SERVER << EOF
    # 创建用户（如果不存在）
    if ! id "todolist" &>/dev/null; then
        sudo useradd -r -s /bin/false todolist
    fi
    
    # 设置权限
    sudo chown -R todolist:todolist $DEPLOY_PATH
    sudo chmod +x $DEPLOY_PATH/bin/$BINARY
    
    # 安装 systemd 服务
    if [ -f "/tmp/$SERVICE_FILE" ]; then
        sudo cp /tmp/$SERVICE_FILE /etc/systemd/system/
        sudo systemctl daemon-reload
        sudo systemctl enable todolist
        echo "systemd 服务已安装"
    fi
    
    # 清理临时文件
    rm -f /tmp/$SERVICE_FILE
    
    echo "部署完成!"
    echo "启动服务: sudo systemctl start todolist"
    echo "查看状态: sudo systemctl status todolist"
    echo "查看日志: sudo journalctl -u todolist -f"
EOF

echo "部署成功!"
echo ""
echo "下一步操作:"
echo "1. 登录服务器: ssh $SERVER"
echo "2. 配置数据库连接等环境变量"
echo "3. 启动服务: sudo systemctl start todolist"
echo "4. 检查状态: sudo systemctl status todolist"
