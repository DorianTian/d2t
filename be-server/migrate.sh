#!/bin/bash

# 确保目录存在
mkdir -p cmd/api
mkdir -p internal/{config,api,handlers,middleware,models,services}

# 移动核心文件
echo "移动核心功能..."
cp -r core internal/core 2>/dev/null || :
cp -r utils internal/utils 2>/dev/null || :

# 更新导入路径
echo "更新导入路径..."
find internal -type f -name "*.go" -exec sed -i '' 's/"d2t_server\/core"/"d2t_server\/internal\/core"/g' {} \;
find internal -type f -name "*.go" -exec sed -i '' 's/"d2t_server\/utils"/"d2t_server\/internal\/utils"/g' {} \;

echo "完成迁移设置！"
echo "请运行 go run cmd/api/main.go 启动新的服务器结构" 