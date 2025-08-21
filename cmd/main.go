package main

import (
	"fmt"
	"os"
	"strings"
	
	"github.com/ct-zh/englishLearn/config"
	"github.com/ct-zh/englishLearn/internal/cli"
	"github.com/ct-zh/englishLearn/internal/dao"
	"github.com/ct-zh/englishLearn/internal/logic/sections"
)

func main() {
	// 获取命令行参数（去除程序名）
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}
	
	// 分离配置参数和应用参数
	configArgs, appArgs := separateArgs(args)
	
	// 创建CLI应用实例
	app, err := createApp(configArgs)
	if err != nil {
		fmt.Printf("初始化应用失败: %v\n", err)
		os.Exit(1)
	}
	
	// 运行应用，传入应用相关的参数
	if err := app.Run(appArgs); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
}

// separateArgs 分离配置参数和应用参数
func separateArgs(args []string) (configArgs []string, appArgs []string) {
	configFlags := map[string]bool{
		"-f": true, "--file": true,
		"-h": true, "--help": true,
	}
	
	i := 0
	for i < len(args) {
		arg := args[i]
		
		// 检查是否是配置相关的参数
		if configFlags[arg] {
			configArgs = append(configArgs, arg)
			// 如果是需要值的参数（-f, --file），也包含下一个参数
			if (arg == "-f" || arg == "--file") && i+1 < len(args) {
				i++
				configArgs = append(configArgs, args[i])
			}
		} else if strings.HasPrefix(arg, "-f=") || strings.HasPrefix(arg, "--file=") {
			// 处理 -f=file.json 或 --file=file.json 格式
			configArgs = append(configArgs, arg)
		} else {
			// 其他参数作为应用参数
			appArgs = append(appArgs, arg)
		}
		i++
	}
	
	return configArgs, appArgs
}

// createApp 创建CLI应用实例
func createApp(configArgs []string) (*cli.App, error) {
	// 加载配置
	cfg, err := config.LoadConfigWithArgs(configArgs)
	if err != nil {
		return nil, err
	}
	
	// 创建DAO工厂
	daoFactory := dao.ProvideDAOFactory(cfg)
	
	// 创建章节DAO
	sectionDAO := dao.ProvideSectionDAO(daoFactory)
	
	// 创建业务逻辑服务
	service := sections.ProvideService(sectionDAO)
	
	// 创建CLI应用
	app := cli.ProvideApp(cfg, service, daoFactory)
	return app, nil
}