package sections

import (
	"fmt"

	"github.com/ct-zh/englishLearn/internal/logic/sections"
	"github.com/ct-zh/englishLearn/model"
)

// CreateSectionNode 创建章节节点
type CreateSectionNode struct {
	*model.BaseMenuNode
	service *sections.Service
}

// NewCreateSection 创建新章节节点
func NewCreateSection(service *sections.Service) *CreateSectionNode {
	node := &CreateSectionNode{
		BaseMenuNode: &model.BaseMenuNode{
			ID:       "createSection",
			Name:     "创建新章节",
			Command:  "1",
			Children: make(map[string]model.MenuNode),
		},
		service: service,
	}

	node.Handler = node.handleCreateSection
	return node
}

// handleCreateSection 处理创建章节的逻辑
func (n *CreateSectionNode) handleCreateSection(ctx *model.MenuContext) error {
	for {
		fmt.Print("请输入新章节名称: ")
		var sectionName string
		if _, err := fmt.Scanln(&sectionName); err != nil {
			fmt.Printf("输入错误: %v\n", err)
			continue
		}

		// 检查输入是否为空
		if sectionName == "" {
			fmt.Println("章节名称不能为空，请重新输入")
			continue
		}

		// 创建章节请求
		req := &model.CreateSectionRequest{
			Name: sectionName,
		}

		// 调用service创建章节
		err := n.service.CreateSection(req)
		if err != nil {
			// 如果是章节已存在的错误，允许用户重新输入
			if fmt.Sprintf("%v", err) == fmt.Sprintf("章节 '%s' 已存在", sectionName) {
				fmt.Printf("错误: %v\n", err)
				fmt.Println("请输入不同的章节名称")
				continue
			}
			// 其他错误直接返回
			return fmt.Errorf("创建章节失败: %w", err)
		}

		// 创建成功，自动选择该章节并进入章节操作菜单
		fmt.Printf("\n章节 '%s' 创建成功！\n", sectionName)

		// 选择刚创建的章节
		selectReq := &model.SelectSectionRequest{
			SectionName: sectionName,
		}

		selectResp, err := n.service.SelectSection(selectReq)
		if err != nil {
			return fmt.Errorf("选择新创建的章节失败: %w", err)
		}

		if selectResp.IsSuccess {
			fmt.Printf("已自动选择章节: %s\n", selectResp.Selected.Name)

			// 创建一个临时的SelectSectionNode来复用章节操作菜单逻辑
			selectNode := NewSelectSection(n.service)
			return selectNode.showSectionMenu(ctx, &selectResp.Selected)
		}

		return nil
	}
}