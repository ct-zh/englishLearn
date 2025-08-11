package model

import (
	"fmt"
)

// MenuNode 菜单节点接口
type MenuNode interface {
	GetID() string
	GetName() string
	GetCommand() string
	GetChildren() map[string]MenuNode
	Menu(child MenuNode) MenuNode // 挂载子节点
	Execute(ctx *MenuContext) error
	Display() string
	IsLeaf() bool
}

// MenuContext 菜单执行上下文
type MenuContext struct {
	CurrentNode MenuNode
	Path        []string                 // 当前路径
	Session     map[string]interface{}   // 会话数据
	Args        map[string]interface{}   // 命令参数
}

// BaseMenuNode 基础菜单节点实现
type BaseMenuNode struct {
	ID       string
	Name     string
	Command  string
	Children map[string]MenuNode
	Handler  func(ctx *MenuContext) error
}

// GetID 获取节点ID
func (b *BaseMenuNode) GetID() string {
	return b.ID
}

// GetName 获取节点名称
func (b *BaseMenuNode) GetName() string {
	return b.Name
}

// GetCommand 获取节点命令
func (b *BaseMenuNode) GetCommand() string {
	return b.Command
}

// GetChildren 获取子节点
func (b *BaseMenuNode) GetChildren() map[string]MenuNode {
	if b.Children == nil {
		b.Children = make(map[string]MenuNode)
	}
	return b.Children
}

// Menu 挂载子节点方法
func (b *BaseMenuNode) Menu(child MenuNode) MenuNode {
	if b.Children == nil {
		b.Children = make(map[string]MenuNode)
	}
	b.Children[child.GetCommand()] = child
	return b
}

// Execute 执行节点处理函数
func (b *BaseMenuNode) Execute(ctx *MenuContext) error {
	if b.Handler != nil {
		return b.Handler(ctx)
	}
	return nil
}

// Display 显示节点信息
func (b *BaseMenuNode) Display() string {
	return fmt.Sprintf("%s. %s", b.Command, b.Name)
}

// IsLeaf 判断是否为叶子节点
func (b *BaseMenuNode) IsLeaf() bool {
	return len(b.Children) == 0
}

// 错误定义
var (
	ErrExit = fmt.Errorf("退出程序")
	ErrBack = fmt.Errorf("返回上级")
)