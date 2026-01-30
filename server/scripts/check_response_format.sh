#!/bin/bash
# 检测 Handler 中是否有直接使用 gin.H 返回响应的代码
# 期望结果：0（表示所有响应都使用了统一响应函数）

HANDLER_DIR="$(dirname "$0")/../internal/handler"

echo "检查 API 响应格式规范..."
echo "================================"

count=$(grep -rn "c.JSON.*gin.H" "$HANDLER_DIR"/*.go 2>/dev/null | grep -v "_test.go" | wc -l | tr -d ' ')

if [ "$count" -eq 0 ]; then
    echo "✅ 通过：所有 Handler 都使用了统一响应函数"
    exit 0
else
    echo "❌ 发现 $count 处违规代码："
    grep -rn "c.JSON.*gin.H" "$HANDLER_DIR"/*.go | grep -v "_test.go"
    echo ""
    echo "请使用 pkg/response 包的统一响应函数替换上述代码"
    exit 1
fi
