#!/bin/bash

# 检查是否有正在运行的 htst 进程
PID=$(pgrep -f "./htst")
if [ -n "$PID" ]; then
    echo "发现正在运行的 htst 进程 (PID: $PID)，正在停止..."
    kill "$PID"
    # 等待进程完全停止
    sleep 3
    # 检查进程是否仍然存在
    if ps -p "$PID" > /dev/null; then
        echo "强制终止进程..."
        kill -9 "$PID" 2>/dev/null
    fi
    echo "旧进程已停止"
fi

# 后台运行新的 htst 进程
nohup ./htst $RUN_MODE > htst.log 2>&1 &
NEW_PID=$!
echo "htst 已在后台启动，PID: $NEW_PID"
echo "日志文件: htst.log"
