package sections

import (
	"fmt"

	"github.com/ct-zh/englishLearn/internal/dao"
	"github.com/ct-zh/englishLearn/internal/logic/sections"
	"github.com/ct-zh/englishLearn/model"
)

// SectionsNode 章节节点
type SectionsNode struct {
	*model.BaseMenuNode
	service *sections.Service
}

// NewSections 创建章节节点
func NewSections() *SectionsNode {
	// 创建DAO工厂和service
	daoFactory := dao.NewDAOFactory("../../data")
	sectionDAO := daoFactory.GetSectionDAO()
	service := sections.NewService(sectionDAO)

	return &SectionsNode{
		BaseMenuNode: &model.BaseMenuNode{
			ID:       "sections",
			Name:     "按章节记忆",
			Command:  "1",
			Children: make(map[string]model.MenuNode),
			Handler: func(ctx *model.MenuContext) error {
				fmt.Println("进入章节管理模式...")
				return nil
			},
		},
		service: service,
	}
}