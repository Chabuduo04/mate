#!/bin/bash

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "Node.js 未安装，请先安装 Node.js 18+"
    exit 1
fi

# 检查npm是否安装
if ! command -v npm &> /dev/null; then
    echo "npm 未安装，请先安装 npm"
    exit 1
fi

# 安装依赖
echo "正在安装依赖..."
npm install

# 启动开发服务器
echo "启动开发服务器..."
npm run dev
