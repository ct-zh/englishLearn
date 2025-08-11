package commands

import (
	"github.com/ct-zh/englishLearn/internal/cli/commands/sections"
	"github.com/ct-zh/englishLearn/model"
)

// MenuRouter 菜单路由器
type MenuRouter struct {
	root model.MenuNode
}

// NewMenuRouter 创建菜单路由器
func NewMenuRouter() *MenuRouter {
	return &MenuRouter{}
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

	// 创建sections节点并挂载到根节点
	sectionsNode := sections.NewSections()
	root.Menu(sectionsNode)

	// 创建selectSection节点并挂载到sections下
	selectSection := sections.NewSelectSection()
	sectionsNode.Menu(selectSection)

	// 创建单词操作节点并挂载到selectSection下
	selectSection.Menu(sections.NewAddWord())
	selectSection.Menu(sections.NewListWords())
	selectSection.Menu(sections.NewRandomWords())

	r.root = root
	return root
}

// GetRoot 获取根节点
func (r *MenuRouter) GetRoot() model.MenuNode {
	return r.root
}