//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/ct-zh/englishLearn/config"
	"github.com/ct-zh/englishLearn/internal/cli"
	"github.com/ct-zh/englishLearn/internal/dao"
	"github.com/ct-zh/englishLearn/internal/logic/sections"
)

// wireApp 使用Wire进行依赖注入，构建完整的应用
func wireApp() (*cli.App, error) {
	wire.Build(
		// Config层
		config.ProvideConfig,
		
		// DAO层
		dao.ProvideDAOFactory,
		dao.ProvideSectionDAO,
		
		// Logic层
		sections.ProvideService,
		
		// CLI层
		cli.ProvideApp,
	)
	return &cli.App{}, nil
}