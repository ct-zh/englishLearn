package cli

import (
	"fmt"
)

// App CLI应用结构体
type App struct {
	// TODO: 添加依赖注入的字段
}

// NewApp 创建新的CLI应用实例
func NewApp() *App {
	return &App{}
}

// Run 运行CLI应用
func (a *App) Run(args []string) error {
	if len(args) == 0 {
		return a.runInteractiveMode()
	}
	return a.runCommandMode(args)
}

// runInteractiveMode 运行交互模式
func (a *App) runInteractiveMode() error {
	fmt.Println("=== 英语学习工具 ===")
	fmt.Println("请选择操作：")
	fmt.Println("1. 添加单词")
	fmt.Println("2. 查看单词列表")
	fmt.Println("3. 搜索单词")
	fmt.Println("4. 开始复习")
	fmt.Println("5. 删除单词")
	fmt.Println("6. 退出")
	fmt.Print("请输入选项 (1-6): ")
	
	// TODO: 实现用户输入处理
	return nil
}

// runCommandMode 运行命令模式
func (a *App) runCommandMode(args []string) error {
	fmt.Printf("执行命令: %v\n", args)
	// TODO: 实现命令解析和执行
	return nil
}