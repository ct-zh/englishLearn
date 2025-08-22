package sections

import (
	"github.com/ct-zh/englishLearn/internal/logic/sections"
	"github.com/ct-zh/englishLearn/model"
)

// ListWordsNode 查看单词节点
type ListWordsNode struct {
	*model.BaseMenuNode
	service *sections.Service
}

// NewListWords 创建查看单词节点
func NewListWords(service *sections.Service) *ListWordsNode {
	return &ListWordsNode{
		BaseMenuNode: &model.BaseMenuNode{
			ID:       "listWords",
			Name:     "查看单词",
			Command:  "3",
			Children: make(map[string]model.MenuNode),
			Handler: func(ctx *model.MenuContext) error {
				req := &model.ListWordsRequest{
					Page: 1,
					Size: 10,
				}

				_, err := service.ListWords(req)
				return err
			},
		},
		service: service,
	}
}