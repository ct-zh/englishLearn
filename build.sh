#!/bin/bash

# 英语学习工具构建脚本
# 编译Go程序到build目录
# 支持指定平台构建

set -e  # 遇到错误时退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目信息
APP_NAME="englishLearn"
BUILD_DIR="build"
MAIN_FILE="./cmd"

# 支持的平台列表
SUPPORTED_PLATFORMS=("darwin-amd64" "darwin-arm64" "linux-amd64" "windows-amd64" "all")

# 自动检测当前平台
detect_current_platform() {
    local os_type arch_type
    
    # 检测操作系统
    os_type=$(uname -s)
    # 检测架构
    arch_type=$(uname -m)
    
    case "$os_type" in
        "Darwin")
            case "$arch_type" in
                "x86_64")
                    echo "darwin-amd64"
                    ;;
                "arm64")
                    echo "darwin-arm64"
                    ;;
                *)
                    echo -e "${YELLOW}警告: 未知的macOS架构 $arch_type，默认使用 darwin-amd64${NC}" >&2
                    echo "darwin-amd64"
                    ;;
            esac
            ;;
        "Linux")
            case "$arch_type" in
                "x86_64")
                    echo "linux-amd64"
                    ;;
                *)
                    echo -e "${YELLOW}警告: 未知的Linux架构 $arch_type，默认使用 linux-amd64${NC}" >&2
                    echo "linux-amd64"
                    ;;
            esac
            ;;
        *)
            echo -e "${YELLOW}警告: 未知的操作系统 $os_type，默认使用 linux-amd64${NC}" >&2
            echo "linux-amd64"
            ;;
    esac
}

# 显示帮助信息
show_help() {
    echo -e "${BLUE}英语学习工具构建脚本${NC}"
    echo
    echo -e "${YELLOW}使用方法:${NC}"
    echo -e "  $0                     构建当前平台 (自动检测)"
    echo -e "  $0 --all               构建所有平台"
    echo -e "  $0 [平台]              构建指定平台"
    echo -e "  $0 --help              显示此帮助信息"
    echo
    echo -e "${YELLOW}支持的平台:${NC}"
    echo -e "  darwin-amd64           macOS (Intel)"
    echo -e "  darwin-arm64           macOS (Apple Silicon)"
    echo -e "  linux-amd64            Linux (x86_64)"
    echo -e "  windows-amd64          Windows (x86_64)"
    echo -e "  all                    所有平台"
    echo
    echo -e "${YELLOW}示例:${NC}"
    echo -e "  $0                     构建当前平台"
    echo -e "  $0 --all               构建所有平台"
    echo -e "  $0 darwin-amd64        只构建macOS Intel版本"
    echo -e "  $0 linux-amd64         只构建Linux版本"
    echo
    echo -e "${YELLOW}当前检测到的平台:${NC} $(detect_current_platform)"
    echo
}

# 验证平台参数
validate_platform() {
    local platform=$1
    for supported in "${SUPPORTED_PLATFORMS[@]}"; do
        if [[ "$platform" == "$supported" ]]; then
            return 0
        fi
    done
    return 1
}

# 解析命令行参数
TARGET_PLATFORM=""
if [[ $# -gt 0 ]]; then
    case $1 in
        --help|-h)
            show_help
            exit 0
            ;;
        --all)
            TARGET_PLATFORM="all"
            ;;
        *)
            if validate_platform "$1"; then
                TARGET_PLATFORM="$1"
            else
                echo -e "${RED}错误: 不支持的平台 '$1'${NC}"
                echo -e "${YELLOW}支持的平台: ${SUPPORTED_PLATFORMS[*]}${NC}"
                echo -e "使用 '$0 --help' 查看详细帮助"
                exit 1
            fi
            ;;
    esac
else
    # 无参数时自动检测当前平台
    TARGET_PLATFORM=$(detect_current_platform)
fi

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: 未找到Go环境，请先安装Go${NC}"
    exit 1
fi

# 显示构建信息
if [[ "$TARGET_PLATFORM" == "all" ]]; then
    echo -e "${YELLOW}开始构建 ${APP_NAME} (所有平台)...${NC}"
else
    echo -e "${YELLOW}开始构建 ${APP_NAME} (${TARGET_PLATFORM})...${NC}"
    if [[ "$TARGET_PLATFORM" == "$(detect_current_platform)" ]]; then
        echo -e "${GREEN}✓ 已自动检测到当前平台${NC}"
    fi
fi

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

# 构建指定平台的函数
build_platform() {
    local platform=$1
    local goos arch ext filename
    
    case $platform in
        "darwin-amd64")
            goos="darwin"
            arch="amd64"
            ext=""
            filename="${APP_NAME}-darwin-amd64"
            echo -e "${YELLOW}编译 macOS (Intel) 版本...${NC}"
            ;;
        "darwin-arm64")
            goos="darwin"
            arch="arm64"
            ext=""
            filename="${APP_NAME}-darwin-arm64"
            echo -e "${YELLOW}编译 macOS (Apple Silicon) 版本...${NC}"
            ;;
        "linux-amd64")
            goos="linux"
            arch="amd64"
            ext=""
            filename="${APP_NAME}-linux-amd64"
            echo -e "${YELLOW}编译 Linux 版本...${NC}"
            ;;
        "windows-amd64")
            goos="windows"
            arch="amd64"
            ext=".exe"
            filename="${APP_NAME}-windows-amd64.exe"
            echo -e "${YELLOW}编译 Windows 版本...${NC}"
            ;;
        *)
            echo -e "${RED}错误: 未知平台 $platform${NC}"
            return 1
            ;;
    esac
    
    GOOS=$goos GOARCH=$arch go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/$filename" "$MAIN_FILE"
    echo -e "${GREEN}✓ $platform 版本编译完成${NC}"
}

# 编译指定平台的二进制文件
echo -e "${YELLOW}编译二进制文件...${NC}"

if [[ "$TARGET_PLATFORM" == "all" ]]; then
    # 构建所有平台
    build_platform "darwin-amd64"
    build_platform "darwin-arm64"
    build_platform "linux-amd64"
    build_platform "windows-amd64"
else
    # 构建指定平台
    build_platform "$TARGET_PLATFORM"
fi

# 创建当前平台的符号链接（仅在构建当前平台或所有平台时）
create_symlink() {
    local current_platform
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        if [[ $(uname -m) == "arm64" ]]; then
            current_platform="darwin-arm64"
        else
            current_platform="darwin-amd64"
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        current_platform="linux-amd64"
    else
        return 0  # 不支持的平台，跳过符号链接创建
    fi
    
    # 检查是否构建了当前平台
    local target_file="$BUILD_DIR/${APP_NAME}-${current_platform}"
    if [[ -f "$target_file" ]]; then
        echo -e "${YELLOW}创建当前平台符号链接...${NC}"
        ln -sf "${APP_NAME}-${current_platform}" "$BUILD_DIR/${APP_NAME}"
        echo -e "${GREEN}✓ 符号链接创建完成: ${APP_NAME} -> ${APP_NAME}-${current_platform}${NC}"
    fi
}

create_symlink

echo
echo -e "${GREEN}构建完成！${NC}"
echo -e "${YELLOW}构建文件位置:${NC}"
if [[ "$TARGET_PLATFORM" == "all" ]]; then
    ls -la "$BUILD_DIR/"
else
    ls -la "$BUILD_DIR/" | grep -E "(total|${TARGET_PLATFORM}|${APP_NAME}$)" || ls -la "$BUILD_DIR/"
fi

echo
echo -e "${YELLOW}使用方法:${NC}"
if [[ -f "$BUILD_DIR/${APP_NAME}" ]]; then
    echo -e "  直接运行: ./$BUILD_DIR/${APP_NAME}"
fi
if [[ "$TARGET_PLATFORM" == "all" ]]; then
    echo -e "  或者指定平台: ./$BUILD_DIR/${APP_NAME}-<platform>"
else
    echo -e "  运行构建的版本: ./$BUILD_DIR/${APP_NAME}-${TARGET_PLATFORM}"
fi
echo
echo -e "${GREEN}构建脚本执行完成！${NC}"