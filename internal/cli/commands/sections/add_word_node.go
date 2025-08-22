package sections

import (
	"fmt"

	"github.com/ct-zh/englishLearn/internal/logic/sections"
	"github.com/ct-zh/englishLearn/model"
)

// AddWordNode 添加单词节点
type AddWordNode struct {
	*model.BaseMenuNode
	service *sections.Service
}

// NewAddWord 创建添加单词节点
func NewAddWord(service *sections.Service) *AddWordNode {
	return &AddWordNode{
		BaseMenuNode: &model.BaseMenuNode{
			ID:       "addWord",
			Name:     "添加单词",
			Command:  "a",
			Children: make(map[string]model.MenuNode),
			Handler: func(ctx *model.MenuContext) error {
				fmt.Println("开始添加单词...")

				// 从上下文参数中获取单词信息
				req := &model.AddWordRequest{
					Section: service.GetCurrentSection(),
				}

				if ctx.Args != nil {
					if word, ok := ctx.Args["word"].(string); ok {
						req.Word = word
					}
					if translation, ok := ctx.Args["translation"].(string); ok {
						req.Translation = translation
					}
				}

				return service.AddWord(req)
			},
		},
		service: service,
	}
}