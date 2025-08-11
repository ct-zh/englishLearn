package main

import (
	"fmt"
	"os"

	"github.com/ct-zh/englishLearn/internal/cli"
)

func main() {
	// 创建CLI应用实例
	app := cli.NewApp()
	
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