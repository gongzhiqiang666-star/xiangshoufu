#!/bin/bash

# 享收付 - 重启所有服务脚本
# 功能：杀掉后端、PC端、APP端，然后重新启动

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="/Users/apple/claudelife/xiangshoufu"

# 日志目录
LOG_DIR="$PROJECT_ROOT/logs"
mkdir -p "$LOG_DIR"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}    享收付 - 重启所有服务${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# ============================================================
# 第一步：杀掉所有相关进程
# ============================================================
echo -e "${YELLOW}[1/2] 停止所有服务...${NC}"

# 杀掉后端进程 (Go server)
echo -n "  - 停止后端服务... "
pkill -f "go run cmd/server/main.go" 2>/dev/null || true
pkill -f "xiangshoufu/server" 2>/dev/null || true
# 杀掉监听8080端口的进程
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
echo -e "${GREEN}完成${NC}"

# 杀掉PC端进程 (Vite dev server)
echo -n "  - 停止PC端服务... "
pkill -f "vite" 2>/dev/null || true
pkill -f "npm run dev" 2>/dev/null || true
# 杀掉监听5173端口的进程
lsof -ti:5173 | xargs kill -9 2>/dev/null || true
echo -e "${GREEN}完成${NC}"

# 杀掉APP端进程 (Flutter)
echo -n "  - 停止APP端服务... "
pkill -f "flutter run" 2>/dev/null || true
pkill -f "flutter_tools" 2>/dev/null || true
echo -e "${GREEN}完成${NC}"

# 等待进程完全退出
sleep 2

echo ""
echo -e "${YELLOW}[2/2] 启动所有服务...${NC}"

# ============================================================
# 第二步：启动所有服务
# ============================================================

# 启动后端服务
echo -n "  - 启动后端服务 (8080)... "
cd "$PROJECT_ROOT/server"
nohup go run cmd/server/main.go > "$LOG_DIR/server.log" 2>&1 &
SERVER_PID=$!
echo -e "${GREEN}PID: $SERVER_PID${NC}"

# 启动APP端服务（先于PC端启动，使用Chrome网页版）
echo -n "  - 启动APP端服务 (Chrome)... "
cd "$PROJECT_ROOT/mobileapp"
nohup flutter run -d chrome > "$LOG_DIR/app.log" 2>&1 &
APP_PID=$!
echo -e "${GREEN}PID: $APP_PID${NC}"

# 启动PC端服务
echo -n "  - 启动PC端服务 (5173)... "
cd "$PROJECT_ROOT/web"
nohup npm run dev > "$LOG_DIR/web.log" 2>&1 &
WEB_PID=$!
echo -e "${GREEN}PID: $WEB_PID${NC}"

# ============================================================
# 第三步：等待并检查服务状态
# ============================================================
echo ""
echo -e "${YELLOW}[检查] 等待服务启动 (最多30秒)...${NC}"

# 检查后端服务
echo -n "  - 后端服务 (8080): "
for i in {1..30}; do
    if curl -s http://localhost:8080/health > /dev/null 2>&1 || curl -s http://localhost:8080/api/health > /dev/null 2>&1 || lsof -ti:8080 > /dev/null 2>&1; then
        echo -e "${GREEN}✓ 启动成功${NC}"
        SERVER_OK=1
        break
    fi
    sleep 1
done
if [ -z "$SERVER_OK" ]; then
    echo -e "${RED}✗ 启动失败${NC}"
    echo -e "${RED}    查看日志: tail -f $LOG_DIR/server.log${NC}"
fi

# 检查PC端服务
echo -n "  - PC端服务 (5173): "
for i in {1..30}; do
    if curl -s http://localhost:5173 > /dev/null 2>&1 || lsof -ti:5173 > /dev/null 2>&1; then
        echo -e "${GREEN}✓ 启动成功${NC}"
        WEB_OK=1
        break
    fi
    sleep 1
done
if [ -z "$WEB_OK" ]; then
    echo -e "${RED}✗ 启动失败${NC}"
    echo -e "${RED}    查看日志: tail -f $LOG_DIR/web.log${NC}"
fi

# 检查APP端服务（Flutter比较特殊，检查进程是否存在）
echo -n "  - APP端服务: "
sleep 5
if ps -p $APP_PID > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 进程运行中 (需要连接设备)${NC}"
    APP_OK=1
else
    echo -e "${YELLOW}⚠ 进程已退出 (可能需要手动启动flutter run)${NC}"
fi

# ============================================================
# 最终结果
# ============================================================
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}    启动结果汇总${NC}"
echo -e "${BLUE}========================================${NC}"

if [ -n "$SERVER_OK" ]; then
    echo -e "  后端服务:  ${GREEN}✓ 运行中${NC}  http://localhost:8080"
else
    echo -e "  后端服务:  ${RED}✗ 失败${NC}"
fi

if [ -n "$WEB_OK" ]; then
    echo -e "  PC端服务:  ${GREEN}✓ 运行中${NC}  http://localhost:5173"
else
    echo -e "  PC端服务:  ${RED}✗ 失败${NC}"
fi

if [ -n "$APP_OK" ]; then
    echo -e "  APP端服务: ${GREEN}✓ 运行中${NC}  (Flutter)"
else
    echo -e "  APP端服务: ${YELLOW}⚠ 需手动启动${NC}"
fi

echo ""
echo -e "${BLUE}日志文件:${NC}"
echo "  - 后端: $LOG_DIR/server.log"
echo "  - PC端: $LOG_DIR/web.log"
echo "  - APP:  $LOG_DIR/app.log"
echo ""

# 返回状态
if [ -n "$SERVER_OK" ] && [ -n "$WEB_OK" ]; then
    echo -e "${GREEN}主要服务启动成功！${NC}"

    # 打开浏览器
    echo ""
    echo -e "${YELLOW}[打开浏览器]${NC}"
    echo -n "  - 打开PC端... "
    open "http://localhost:5173" 2>/dev/null && echo -e "${GREEN}完成${NC}" || echo -e "${RED}失败${NC}"

    sleep 1

    echo -n "  - 打开APP端... "
    open "http://localhost:8080" 2>/dev/null && echo -e "${GREEN}完成${NC}" || echo -e "${RED}失败${NC}"
    echo ""

    exit 0
else
    echo -e "${RED}部分服务启动失败，请检查日志${NC}"
    exit 1
fi
