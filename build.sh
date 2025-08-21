#!/bin/bash

# 英语学习工具构建脚本
# 编译Go程序到build目录

set -e  # 遇到错误时退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目信息
APP_NAME="englishLearn"
BUILD_DIR="build"
MAIN_FILE="./cmd"

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: 未找到Go环境，请先安装Go${NC}"
    exit 1
fi

echo -e "${YELLOW}开始构建 ${APP_NAME}...${NC}"

# 创建build目录（如果不存在）
if [ ! -d "$BUILD_DIR" ]; then
    echo -e "${YELLOW}创建build目录...${NC}"
    mkdir -p "$BUILD_DIR"
fi

# 获取版本信息
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建标志
LDFLAGS="-X 'main.Version=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.GitCommit=${GIT_COMMIT}'"

echo -e "${YELLOW}版本信息:${NC}"
echo -e "  版本: ${VERSION}"
echo -e "  构建时间: ${BUILD_TIME}"
echo -e "  Git提交: ${GIT_COMMIT}"
echo

# 编译不同平台的二进制文件
echo -e "${YELLOW}编译二进制文件...${NC}"

# macOS (当前平台)
echo -e "${YELLOW}编译 macOS 版本...${NC}"
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/${APP_NAME}-darwin-amd64" "$MAIN_FILE"
echo -e "${GREEN}✓ macOS (Intel) 版本编译完成${NC}"

GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/${APP_NAME}-darwin-arm64" "$MAIN_FILE"
echo -e "${GREEN}✓ macOS (Apple Silicon) 版本编译完成${NC}"

# Linux
echo -e "${YELLOW}编译 Linux 版本...${NC}"
GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/${APP_NAME}-linux-amd64" "$MAIN_FILE"
echo -e "${GREEN}✓ Linux 版本编译完成${NC}"

# Windows
echo -e "${YELLOW}编译 Windows 版本...${NC}"
GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/${APP_NAME}-windows-amd64.exe" "$MAIN_FILE"
echo -e "${GREEN}✓ Windows 版本编译完成${NC}"

# 创建当前平台的符号链接
echo -e "${YELLOW}创建当前平台符号链接...${NC}"
if [[ "$OSTYPE" == "darwin"* ]]; then
    if [[ $(uname -m) == "arm64" ]]; then
        ln -sf "${APP_NAME}-darwin-arm64" "$BUILD_DIR/${APP_NAME}"
    else
        ln -sf "${APP_NAME}-darwin-amd64" "$BUILD_DIR/${APP_NAME}"
    fi
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    ln -sf "${APP_NAME}-linux-amd64" "$BUILD_DIR/${APP_NAME}"
fi

echo
echo -e "${GREEN}构建完成！${NC}"
echo -e "${YELLOW}构建文件位置:${NC}"
ls -la "$BUILD_DIR/"

echo
echo -e "${YELLOW}使用方法:${NC}"
echo -e "  直接运行: ./$BUILD_DIR/${APP_NAME}"
echo -e "  或者: ./$BUILD_DIR/${APP_NAME}-<platform>-<arch>"
echo
echo -e "${GREEN}构建脚本执行完成！${NC}"