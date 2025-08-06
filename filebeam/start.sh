#!/bin/bash

echo "========================================"
echo "FileBeam 文件共享服务启动脚本"
echo "========================================"
echo

# 检查是否存在可执行文件
if [ ! -f "filebeam" ]; then
    echo "正在编译项目..."
    go build -o filebeam main.go
    if [ $? -ne 0 ]; then
        echo "编译失败！请检查Go环境是否正确安装。"
        exit 1
    fi
fi

echo "启动 FileBeam 服务..."
echo
echo "提示："
echo "- 程序启动后会询问您输入共享文件夹路径"
echo "- 默认访问地址：http://localhost:8888/"
echo "- 默认上传密码：123456"
echo "- 按 Ctrl+C 停止服务"
echo
echo "========================================"

./filebeam 