package sections

import (
	"fmt"

	"github.com/ct-zh/englishLearn/internal/logic/sections"
	"github.com/ct-zh/englishLearn/model"
)

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
// 返回值说明:
// - nil: 成功选择章节并进入章节操作菜单
// - model.ErrBack: 用户选择返回上级菜单
// - 其他error: 发生错误
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
		if _, err := fmt.Scanln(&input); err != nil {
			fmt.Printf("输入错误: %v\n", err)
			continue
		}
		
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
		if _, err := fmt.Scanln(&choice); err != nil {
			fmt.Printf("输入错误: %v\n", err)
			continue
		}

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
			// 重新选择章节，如果用户在章节列表中选择返回，则直接返回上级菜单
			if err := n.handleSelectSection(ctx); err != nil {
				if err == model.ErrBack {
					return model.ErrBack
				}
				fmt.Printf("选择章节失败: %v\n", err)
			}
			// 如果成功选择了新章节，会返回新的章节操作菜单，这里不需要额外处理
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
	if _, err := fmt.Scanln(&word); err != nil {
		return fmt.Errorf("输入错误: %v", err)
	}

	fmt.Print("请输入中文释义: ")
	var translation string
	if _, err := fmt.Scanln(&translation); err != nil {
		return fmt.Errorf("输入错误: %v", err)
	}

	fmt.Print("请输入例句(可选，直接回车跳过): ")
	var phrase string
	// 例句是可选的，忽略输入错误
	_, _ = fmt.Scanln(&phrase)

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
	if _, err := fmt.Scanln(&input); err != nil {
		// 输入错误时使用默认值
		input = ""
	}

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
	if _, err := fmt.Scanln(&keyword); err != nil {
		return fmt.Errorf("输入错误: %v", err)
	}

	req := &model.SearchWordRequest{
		Keyword: keyword,
		Section: sectionName,
	}

	_, err := n.service.SearchWord(req)
	return err
}