package commands

import (
	"github.com/ct-zh/englishLearn/internal/cli/commands/sections"
	"github.com/ct-zh/englishLearn/internal/dao"
	sectionsLogic "github.com/ct-zh/englishLearn/internal/logic/sections"
	"github.com/ct-zh/englishLearn/model"
)

// MenuRouter 菜单路由器
type MenuRouter struct {
	root       model.MenuNode
	service    *sectionsLogic.Service
	daoFactory *dao.DAOFactory
}

// NewMenuRouter 创建菜单路由器
func NewMenuRouter() *MenuRouter {
	return &MenuRouter{}
}

// NewMenuRouterWithService 创建带service的菜单路由器 (用于Wire)
func NewMenuRouterWithService(service *sectionsLogic.Service, daoFactory *dao.DAOFactory) *MenuRouter {
	return &MenuRouter{
		service:    service,
		daoFactory: daoFactory,
	}
}

// newRoot 创建根节点（避免循环导入）
func (r *MenuRouter) newRoot() model.MenuNode {
	return &model.BaseMenuNode{
		ID:       "root",
		Name:     "英语学习工具",
		Command:  "",
		Children: make(map[string]model.MenuNode),
		Handler: func(ctx *model.MenuContext) error {
			// 根节点不执行具体操作，只显示菜单
			return nil
		},
	}
}

// BuildDefaultTree 构建默认菜单树
func (r *MenuRouter) BuildDefaultTree() model.MenuNode {
	// 创建根节点
	root := r.newRoot()

	// 使用注入的service或创建默认service
	service := r.service
	if service == nil {
		// 兼容旧的方式，用于非Wire场景
		daoFactory := dao.NewDAOFactory("../../data")
		sectionDAO := daoFactory.GetSectionDAO()
		service = sectionsLogic.NewService(sectionDAO)
	}

	// 创建sections节点并挂载到根节点
	sectionsNode := sections.NewSections()
	root.Menu(sectionsNode)

	// 创建createSection节点并挂载到sections下（第一个选项）
	createSection := sections.NewCreateSection(service)
	sectionsNode.Menu(createSection)

	// 创建selectSection节点并挂载到sections下
	selectSection := sections.NewSelectSection(service)
	sectionsNode.Menu(selectSection)

	// 创建单词操作节点并挂载到selectSection下
	selectSection.Menu(sections.NewAddWord(service))
	selectSection.Menu(sections.NewListWords(service))
	selectSection.Menu(sections.NewRandomWords(service))

	// 创建文件管理节点并挂载到根节点
	if r.daoFactory != nil {
		fileManager := NewFileManager(r.daoFactory)
		root.Menu(fileManager)
	}

	r.root = root
	return root
}

// GetRoot 获取根节点
func (r *MenuRouter) GetRoot() model.MenuNode {
	return r.root
}