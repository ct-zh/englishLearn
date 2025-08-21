package cli

import (
	"fmt"
	"github.com/ct-zh/englishLearn/internal/cli/commands"
	"github.com/ct-zh/englishLearn/internal/dao"
	sectionsLogic "github.com/ct-zh/englishLearn/internal/logic/sections"
	"github.com/ct-zh/englishLearn/model"
)

// MenuTreeBuilder 菜单树构建器
type MenuTreeBuilder struct {
	root   model.MenuNode
	router *commands.MenuRouter
}

// NewMenuTreeBuilder 创建菜单树构建器
func NewMenuTreeBuilder() *MenuTreeBuilder {
	return &MenuTreeBuilder{
		router: commands.NewMenuRouter(),
	}
}

// NewMenuTreeBuilderWithService 创建带service的菜单树构建器 (用于Wire)
func NewMenuTreeBuilderWithService(service *sectionsLogic.Service, daoFactory *dao.DAOFactory) *MenuTreeBuilder {
	return &MenuTreeBuilder{
		router: commands.NewMenuRouterWithService(service, daoFactory),
	}
}

// BuildDefaultTree 构建默认菜单树
func (b *MenuTreeBuilder) BuildDefaultTree() model.MenuNode {
	// 使用路由器构建菜单树
	root := b.router.BuildDefaultTree()
	b.root = root
	return root
}

// ValidateTree 验证菜单树（检查命令冲突）
func (b *MenuTreeBuilder) ValidateTree(node model.MenuNode) error {
	return b.validateNode(node, []string{})
}

// validateNode 递归验证节点
func (b *MenuTreeBuilder) validateNode(node model.MenuNode, path []string) error {
	currentPath := append(path, node.GetID())
	
	// 检查当前节点下的子节点命令是否冲突
	commands := make(map[string]string)
	for _, child := range node.GetChildren() {
		cmd := child.GetCommand()
		if cmd == "" {
			continue // 跳过空命令（如根节点）
		}
		
		if existingNodeID, exists := commands[cmd]; exists {
			return fmt.Errorf("命令冲突: 命令 '%s' 在节点 '%s' 下被 '%s' 和 '%s' 同时使用", 
				cmd, node.GetID(), existingNodeID, child.GetID())
		}
		commands[cmd] = child.GetID()
	}

	// 递归验证子节点
	for _, child := range node.GetChildren() {
		if err := b.validateNode(child, currentPath); err != nil {
			return err
		}
	}

	return nil
}

// GetRoot 获取根节点
func (b *MenuTreeBuilder) GetRoot() model.MenuNode {
	return b.root
}

// PrintTree 打印菜单树结构（用于调试）
func (b *MenuTreeBuilder) PrintTree(node model.MenuNode, indent string) {
	if node == nil {
		return
	}
	
	fmt.Printf("%s%s (%s) [命令: %s]\n", indent, node.GetName(), node.GetID(), node.GetCommand())
	
	for _, child := range node.GetChildren() {
		b.PrintTree(child, indent+"  ")
	}
}