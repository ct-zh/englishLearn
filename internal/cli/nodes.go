package cli

import (
	"github.com/ct-zh/englishLearn/model"
)

// RootNode 根节点
type RootNode struct {
	*model.BaseMenuNode
}

// NewRoot 创建根节点
func NewRoot() *RootNode {
	return &RootNode{
		BaseMenuNode: &model.BaseMenuNode{
			ID:       "root",
			Name:     "英语学习工具",
			Command:  "",
			Children: make(map[string]model.MenuNode),
			Handler: func(ctx *model.MenuContext) error {
				// 根节点不执行具体操作，只显示菜单
				return nil
			},
		},
	}
}
