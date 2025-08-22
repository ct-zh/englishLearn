package sections

import (
	"github.com/ct-zh/englishLearn/internal/logic/sections"
	"github.com/ct-zh/englishLearn/model"
)

// RandomWordsNode 随机练习节点
type RandomWordsNode struct {
	*model.BaseMenuNode
	service *sections.Service
}

// NewRandomWords 创建随机练习节点
func NewRandomWords(service *sections.Service) *RandomWordsNode {
	return &RandomWordsNode{
		BaseMenuNode: &model.BaseMenuNode{
			ID:       "randomWords",
			Name:     "随机练习",
			Command:  "4",
			Children: make(map[string]model.MenuNode),
			Handler: func(ctx *model.MenuContext) error {
				req := &model.RandomWordsRequest{
					Count: 10,
				}

				_, err := service.RandomWords(req)
				return err
			},
		},
		service: service,
	}
}