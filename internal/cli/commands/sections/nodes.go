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

// SelectSectionNode 选择章节节点
type SelectSectionNode struct {
	*model.BaseMenuNode
	service     *sections.Service
	currentPage int
	pageSize    int
}

// NewSelectSection 创建选择章节节点
func NewSelectSection(service *sections.Service) *SelectSectionNode {
	node := &SelectSectionNode{
		BaseMenuNode: &model.BaseMenuNode{
			ID:       "selectSection",
			Name:     "选择章节",
			Command:  "2",
			Children: make(map[string]model.MenuNode),
		},
		service:     service,
		currentPage: 1,
		pageSize:    5, // 每页显示5个章节
	}
	
	node.Handler = node.handleSelectSection
	return node
}

// handleSelectSection 处理选择章节的逻辑
func (n *SelectSectionNode) handleSelectSection(ctx *model.MenuContext) error {
	for {
		// 获取章节列表
		req := &model.ListSectionsRequest{
			Page: n.currentPage,
			Size: n.pageSize,
		}
		
		resp, err := n.service.ListSections(req)
		if err != nil {
			fmt.Printf("获取章节列表失败: %v\n", err)
			return err
		}
		
		if len(resp.Sections) == 0 {
			fmt.Println("没有找到任何章节")
			return nil
		}
		
		// 显示章节列表
		fmt.Printf("\n=== 章节列表 (第%d页/共%d页) ===\n", resp.CurrentPage, resp.TotalPages)
		for i, section := range resp.Sections {
			fmt.Printf("%d. %s (包含 %d 个单词)\n", i+1, section.Name, len(section.Words))
		}
		
		// 显示操作选项
		fmt.Println("\n操作选项:")
		if resp.HasPrev {
			fmt.Println("p. 上一页")
		}
		if resp.HasNext {
			fmt.Println("n. 下一页")
		}
		fmt.Println("b. 返回上级菜单")
		fmt.Printf("请选择章节序号(1-%d)或操作: ", len(resp.Sections))
		
		// 读取用户输入
		var input string
		fmt.Scanln(&input)
		
		switch input {
		case "p":
			if resp.HasPrev {
				n.currentPage--
			} else {
				fmt.Println("已经是第一页了")
			}
		case "n":
			if resp.HasNext {
				n.currentPage++
			} else {
				fmt.Println("已经是最后一页了")
			}
		case "b":
			return model.ErrBack
		default:
			// 尝试解析为数字
			if choice := parseChoice(input, len(resp.Sections)); choice > 0 {
				selectedSection := resp.Sections[choice-1]
				
				// 选择章节
				selectReq := &model.SelectSectionRequest{
					SectionName: selectedSection.Name,
				}
				
				selectResp, err := n.service.SelectSection(selectReq)
				if err != nil {
					fmt.Printf("选择章节失败: %v\n", err)
					continue
				}
				
				if selectResp.IsSuccess {
					fmt.Printf("\n✓ 已选择章节: %s (包含 %d 个单词)\n", 
						selectResp.Selected.Name, selectResp.WordCount)
					
					// 显示章节操作菜单
					return n.showSectionMenu(ctx, &selectResp.Selected)
				}
			} else {
				fmt.Println("无效的选择，请重新输入")
			}
		}
	}
}

// showSectionMenu 显示章节操作菜单
func (n *SelectSectionNode) showSectionMenu(ctx *model.MenuContext, section *model.SectionEntity) error {
	for {
		fmt.Printf("\n=== 章节: %s ===\n", section.Name)
		fmt.Println("1. 添加单词")
		fmt.Println("2. 查看单词列表")
		fmt.Println("3. 随机练习")
		fmt.Println("4. 搜索单词")
		fmt.Println("5. 重新选择章节")
		fmt.Println("b. 返回上级菜单")
		fmt.Print("请选择操作: ")
		
		var choice string
		fmt.Scanln(&choice)
		
		switch choice {
		case "1":
			if err := n.handleAddWord(section.Name); err != nil {
				fmt.Printf("添加单词失败: %v\n", err)
			}
		case "2":
			if err := n.handleListWords(section.Name); err != nil {
				fmt.Printf("查看单词列表失败: %v\n", err)
			}
		case "3":
			if err := n.handleRandomWords(section.Name); err != nil {
				fmt.Printf("随机练习失败: %v\n", err)
			}
		case "4":
			if err := n.handleSearchWords(section.Name); err != nil {
				fmt.Printf("搜索单词失败: %v\n", err)
			}
		case "5":
			return n.handleSelectSection(ctx) // 重新选择章节
		case "b":
			return model.ErrBack
		default:
			fmt.Println("无效的选择，请重新输入")
		}
	}
}

// handleAddWord 处理添加单词
func (n *SelectSectionNode) handleAddWord(sectionName string) error {
	fmt.Print("请输入单词: ")
	var word string
	fmt.Scanln(&word)
	
	fmt.Print("请输入中文释义: ")
	var translation string
	fmt.Scanln(&translation)
	
	fmt.Print("请输入例句(可选，直接回车跳过): ")
	var phrase string
	fmt.Scanln(&phrase)
	
	req := &model.AddWordRequest{
		Word:        word,
		Translation: translation,
		Phrase:      phrase,
		Section:     sectionName,
	}
	
	return n.service.AddWord(req)
}

// handleListWords 处理查看单词列表
func (n *SelectSectionNode) handleListWords(sectionName string) error {
	req := &model.ListWordsRequest{
		Section: sectionName,
		Page:    1,
		Size:    10,
	}
	
	_, err := n.service.ListWords(req)
	return err
}

// handleRandomWords 处理随机练习
func (n *SelectSectionNode) handleRandomWords(sectionName string) error {
	fmt.Print("请输入练习单词数量(默认10): ")
	var input string
	fmt.Scanln(&input)
	
	count := 10
	if c := parseChoice(input, 100); c > 0 {
		count = c
	}
	
	req := &model.RandomWordsRequest{
		Section: sectionName,
		Count:   count,
	}
	
	_, err := n.service.RandomWords(req)
	return err
}

// handleSearchWords 处理搜索单词
func (n *SelectSectionNode) handleSearchWords(sectionName string) error {
	fmt.Print("请输入搜索关键词: ")
	var keyword string
	fmt.Scanln(&keyword)
	
	req := &model.SearchWordRequest{
		Keyword: keyword,
		Section: sectionName,
	}
	
	_, err := n.service.SearchWord(req)
	return err
}

// parseChoice 解析用户输入的选择
func parseChoice(input string, max int) int {
	if input == "" {
		return 0
	}
	
	// 简单的字符串转数字
	choice := 0
	for _, r := range input {
		if r >= '0' && r <= '9' {
			choice = choice*10 + int(r-'0')
		} else {
			return 0
		}
	}
	
	if choice >= 1 && choice <= max {
		return choice
	}
	return 0
}

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
		fmt.Scanln(&sectionName)
		
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