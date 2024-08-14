#!/bin/bash

# 服务名称
SERVICE_NAME="zhatuche"

# 服务描述
SERVICE_DESCRIPTION="渣土车20240812"

# 当前脚本所在目录
CURRENT_DIR=$(pwd)

# 服务可执行文件路径
SERVICE_EXECUTABLE_PATH="$CURRENT_DIR/run.py"

# 服务配置文件路径
SERVICE_CONFIG_PATH="$CURRENT_DIR/config.conf"

# 服务日志文件路径
SERVICE_LOG_PATH="/var/log/your_service.log"

# 创建 systemd 服务文件
cat << EOF > /etc/systemd/system/$SERVICE_NAME.service
[Unit]
Description=$SERVICE_DESCRIPTION
After=network.target

[Service]
Type=simple
ExecStart=python3 $SERVICE_EXECUTABLE_PATH
WorkingDirectory=$CURRENT_DIR

Restart=always

[Install]
WantedBy=multi-user.target
EOF

# 启用服务
systemctl enable $SERVICE_NAME

# 启动服务
systemctl start $SERVICE_NAME