package cli

import (
	"fmt"

	"github.com/ct-zh/englishLearn/config"
	"github.com/ct-zh/englishLearn/internal/dao"
	sectionsLogic "github.com/ct-zh/englishLearn/internal/logic/sections"
)

// App CLI应用结构
type App struct {
	name       string
	builder    *MenuTreeBuilder
	resolver   *CommandPathResolver
	config     *config.Config
	service    *sectionsLogic.Service
}

// NewApp 创建新的CLI应用
func NewApp() *App {
	builder := NewMenuTreeBuilder()
	root := builder.BuildDefaultTree()
	
	// 验证菜单树
	if err := builder.ValidateTree(root); err != nil {
		panic(fmt.Sprintf("菜单树验证失败: %v", err))
	}
	
	resolver := NewCommandPathResolver(root)
	
	return &App{
		name:     "英语学习工具",
		builder:  builder,
		resolver: resolver,
	}
}

// NewAppWithService 创建带有service的CLI应用 (用于Wire)
func NewAppWithService(cfg *config.Config, service *sectionsLogic.Service, daoFactory *dao.DAOFactory) *App {
	builder := NewMenuTreeBuilderWithService(service, daoFactory)
	root := builder.BuildDefaultTree()
	
	// 验证菜单树
	if err := builder.ValidateTree(root); err != nil {
		panic(fmt.Sprintf("菜单树验证失败: %v", err))
	}
	
	resolver := NewCommandPathResolver(root)
	
	return &App{
		name:     "英语学习工具",
		builder:  builder,
		resolver: resolver,
		config:   cfg,
		service:  service,
	}
}

// ProvideApp 提供CLI应用实例 (Wire Provider)
func ProvideApp(cfg *config.Config, service *sectionsLogic.Service, daoFactory *dao.DAOFactory) *App {
	return NewAppWithService(cfg, service, daoFactory)
}

// Run 运行CLI应用
func (a *App) Run(args []string) error {
	if len(args) > 0 {
		// 命令行模式
		return a.runCommandMode(args)
	} else {
		// 交互模式
		return a.runInteractiveMode()
	}
}

// runCommandMode 运行命令行模式
func (a *App) runCommandMode(args []string) error {
	return a.resolver.ExecuteCommand(args)
}

// runInteractiveMode 运行交互模式
func (a *App) runInteractiveMode() error {
	root := a.builder.GetRoot()
	engine := NewInteractiveEngineWithConfig(root, a.config)
	return engine.Start()
}

// PrintMenuTree 打印菜单树（调试用）
func (a *App) PrintMenuTree() {
	a.builder.PrintTree(a.builder.GetRoot(), "")
}

// ListCommands 列出所有可用命令
func (a *App) ListCommands() {
	a.resolver.ListCommands()
}