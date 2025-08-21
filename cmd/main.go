package main

import (
	"fmt"
	"os"
)

func main() {
	// 使用Wire进行依赖注入，创建CLI应用实例
	app, err := wireApp()
	if err != nil {
		fmt.Printf("初始化应用失败: %v\n", err)
		os.Exit(1)
	}
	
	// 运行应用，传入命令行参数（去除程序名）
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}
	
	if err := app.Run(args); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
}