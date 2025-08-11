package sections

import (
	"fmt"
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
	service := sections.NewService()
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

// SelectSectionNode 选择章节节点
type SelectSectionNode struct {
	*model.BaseMenuNode
	service *sections.Service
}

// NewSelectSection 创建选择章节节点
func NewSelectSection() *SelectSectionNode {
	service := sections.NewService()
	return &SelectSectionNode{
		BaseMenuNode: &model.BaseMenuNode{
			ID:       "selectSection",
			Name:     "选择章节",
			Command:  "1",
			Children: make(map[string]model.MenuNode),
			Handler: func(ctx *model.MenuContext) error {
				currentSection := service.GetCurrentSection()
				fmt.Printf("当前章节: %s\n", currentSection)
				return nil
			},
		},
		service: service,
	}
}

// AddWordNode 添加单词节点
type AddWordNode struct {
	*model.BaseMenuNode
	service *sections.Service
}

// NewAddWord 创建添加单词节点
func NewAddWord() *AddWordNode {
	service := sections.NewService()
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

// ListWordsNode 查看单词节点
type ListWordsNode struct {
	*model.BaseMenuNode
	service *sections.Service
}

// NewListWords 创建查看单词节点
func NewListWords() *ListWordsNode {
	service := sections.NewService()
	return &ListWordsNode{
		BaseMenuNode: &model.BaseMenuNode{
			ID:       "listWords",
			Name:     "查看单词",
			Command:  "l",
			Children: make(map[string]model.MenuNode),
			Handler: func(ctx *model.MenuContext) error {
				fmt.Println("显示单词列表...")
				
				req := &model.ListWordsRequest{
					Section: service.GetCurrentSection(),
					Page:    1,
					Size:    10,
				}
				
				if ctx.Args != nil {
					if page, ok := ctx.Args["page"].(int); ok {
						req.Page = page
					}
					if size, ok := ctx.Args["size"].(int); ok {
						req.Size = size
					}
				}
				
				_, err := service.ListWords(req)
				return err
			},
		},
		service: service,
	}
}

// RandomWordsNode 随机练习节点
type RandomWordsNode struct {
	*model.BaseMenuNode
	service *sections.Service
}

// NewRandomWords 创建随机练习节点
func NewRandomWords() *RandomWordsNode {
	service := sections.NewService()
	return &RandomWordsNode{
		BaseMenuNode: &model.BaseMenuNode{
			ID:       "randomWords",
			Name:     "随机练习",
			Command:  "r",
			Children: make(map[string]model.MenuNode),
			Handler: func(ctx *model.MenuContext) error {
				count := 10 // 默认值
				if ctx.Args != nil {
					if c, ok := ctx.Args["count"].(int); ok {
						count = c
					}
				}
				
				req := &model.RandomWordsRequest{
					Section: service.GetCurrentSection(),
					Count:   count,
				}
				
				fmt.Printf("开始随机练习 %d 个单词...\n", count)
				_, err := service.RandomWords(req)
				return err
			},
		},
		service: service,
	}
}