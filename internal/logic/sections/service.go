package sections

import (
	"context"
	"fmt"
	"math"
	"github.com/ct-zh/englishLearn/internal/dao"
	"github.com/ct-zh/englishLearn/model"
)

// Service sections业务逻辑服务
type Service struct {
	sectionDAO     dao.SectionDAOInterface
	currentSection string // 当前选中的章节
}

// NewService 创建新的sections服务实例
func NewService(sectionDAO dao.SectionDAOInterface) *Service {
	return &Service{
		sectionDAO: sectionDAO,
	}
}

// ProvideService 提供sections服务实例 (Wire Provider)
func ProvideService(sectionDAO dao.SectionDAOInterface) *Service {
	return NewService(sectionDAO)
}

// AddWord 添加单词
func (s *Service) AddWord(req *model.AddWordRequest) error {
	ctx := context.Background()
	
	// 检查章节是否存在
	exists, err := s.sectionDAO.SectionExists(ctx, req.Section)
	if err != nil {
		return fmt.Errorf("检查章节存在性失败: %w", err)
	}
	
	if !exists {
		return fmt.Errorf("章节 '%s' 不存在", req.Section)
	}
	
	// 创建单词实体
	word := model.WordEntity{
		W:      req.Word,
		C:      req.Translation,
		Phrase: "", // 暂时为空，后续可以扩展
	}
	
	// 添加单词到章节
	err = s.sectionDAO.AddWordToSection(ctx, req.Section, word)
	if err != nil {
		return fmt.Errorf("添加单词失败: %w", err)
	}
	
	fmt.Printf("成功添加单词: %s (%s) 到章节: %s\n", req.Word, req.Translation, req.Section)
	return nil
}

// ListWords 获取单词列表
func (s *Service) ListWords(req *model.ListWordsRequest) (*model.ListWordsResponse, error) {
	ctx := context.Background()
	
	// 获取章节
	section, err := s.sectionDAO.GetSection(ctx, req.Section)
	if err != nil {
		return nil, fmt.Errorf("获取章节失败: %w", err)
	}
	
	// 计算分页
	total := len(section.Words)
	totalPages := int(math.Ceil(float64(total) / float64(req.Size)))
	if req.Page > totalPages {
		req.Page = totalPages
	}
	if req.Page < 1 {
		req.Page = 1
	}
	
	start := (req.Page - 1) * req.Size
	end := start + req.Size
	
	if start >= total {
		return &model.ListWordsResponse{
			Words:       []model.WordEntity{},
			Total:       total,
			CurrentPage: req.Page,
			TotalPages:  totalPages,
			HasNext:     false,
			HasPrev:     req.Page > 1,
		}, nil
	}
	
	if end > total {
		end = total
	}
	
	words := section.Words[start:end]
	fmt.Printf("章节 %s 第%d页单词列表 (第%d页/共%d页):\n", req.Section, req.Page, req.Page, totalPages)
	for i, word := range words {
		fmt.Printf("%d. %s - %s\n", start+i+1, word.W, word.C)
	}
	
	return &model.ListWordsResponse{
		Words:       words,
		Total:       total,
		CurrentPage: req.Page,
		TotalPages:  totalPages,
		HasNext:     req.Page < totalPages,
		HasPrev:     req.Page > 1,
	}, nil
}

// RandomWords 随机练习单词
func (s *Service) RandomWords(req *model.RandomWordsRequest) (*model.RandomWordsResponse, error) {
	ctx := context.Background()
	
	// 获取章节
	section, err := s.sectionDAO.GetSection(ctx, req.Section)
	if err != nil {
		return nil, fmt.Errorf("获取章节失败: %w", err)
	}
	
	if len(section.Words) == 0 {
		return &model.RandomWordsResponse{
			Words: []model.WordEntity{},
			Count: 0,
		}, fmt.Errorf("章节 '%s' 中没有单词", req.Section)
	}
	
	// 简单的随机选择（这里可以后续优化为更好的随机算法）
	count := req.Count
	if count > len(section.Words) {
		count = len(section.Words)
	}
	
	// 取前count个单词作为示例（实际应该随机选择）
	randomWords := section.Words[:count]
	
	fmt.Printf("从章节 %s 随机选择 %d 个单词进行练习:\n", req.Section, count)
	for i, word := range randomWords {
		fmt.Printf("%d. %s - %s\n", i+1, word.W, word.C)
	}
	
	return &model.RandomWordsResponse{
		Words: randomWords,
		Count: count,
	}, nil
}

// SearchWord 搜索单词
func (s *Service) SearchWord(req *model.SearchWordRequest) (*model.SearchWordResponse, error) {
	ctx := context.Background()
	
	var searchWords []model.WordEntity
	
	if req.Section != "" {
		// 在指定章节中搜索
		section, err := s.sectionDAO.GetSection(ctx, req.Section)
		if err != nil {
			return nil, fmt.Errorf("获取章节失败: %w", err)
		}
		
		for _, word := range section.Words {
			if contains(word.W, req.Keyword) || contains(word.C, req.Keyword) {
				searchWords = append(searchWords, word)
			}
		}
	} else {
		// 在所有章节中搜索
		allSections, err := s.sectionDAO.ListSections(ctx)
		if err != nil {
			return nil, fmt.Errorf("获取所有章节失败: %w", err)
		}
		
		for _, section := range allSections {
			for _, word := range section.Words {
				if contains(word.W, req.Keyword) || contains(word.C, req.Keyword) {
					searchWords = append(searchWords, word)
				}
			}
		}
	}
	
	fmt.Printf("搜索关键词 '%s' 找到 %d 个结果:\n", req.Keyword, len(searchWords))
	for i, word := range searchWords {
		fmt.Printf("%d. %s - %s\n", i+1, word.W, word.C)
	}
	
	return &model.SearchWordResponse{
		Words: searchWords,
		Total: len(searchWords),
	}, nil
}

// contains 简单的字符串包含检查（不区分大小写）
func contains(str, substr string) bool {
	return len(str) >= len(substr) && 
		   (str == substr || 
		    len(substr) == 0 || 
		    indexOfIgnoreCase(str, substr) >= 0)
}

// indexOfIgnoreCase 不区分大小写的字符串查找
func indexOfIgnoreCase(str, substr string) int {
	strLower := toLower(str)
	substrLower := toLower(substr)
	
	for i := 0; i <= len(strLower)-len(substrLower); i++ {
		if strLower[i:i+len(substrLower)] == substrLower {
			return i
		}
	}
	return -1
}

// toLower 简单的转小写函数
func toLower(s string) string {
	result := make([]byte, len(s))
	for i, b := range []byte(s) {
		if b >= 'A' && b <= 'Z' {
			result[i] = b + 32
		} else {
			result[i] = b
		}
	}
	return string(result)
}

// ListSections 分页获取章节列表
func (s *Service) ListSections(req *model.ListSectionsRequest) (*model.ListSectionsResponse, error) {
	ctx := context.Background()
	
	// 获取所有章节
	allSections, err := s.sectionDAO.ListSections(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取章节列表失败: %w", err)
	}
	
	total := len(allSections)
	if total == 0 {
		return &model.ListSectionsResponse{
			Sections:    []model.SectionEntity{},
			Total:       0,
			CurrentPage: req.Page,
			TotalPages:  0,
			HasNext:     false,
			HasPrev:     false,
		}, nil
	}
	
	// 计算分页
	totalPages := int(math.Ceil(float64(total) / float64(req.Size)))
	if req.Page > totalPages {
		req.Page = totalPages
	}
	if req.Page < 1 {
		req.Page = 1
	}
	
	// 计算起始和结束索引
	start := (req.Page - 1) * req.Size
	end := start + req.Size
	if end > total {
		end = total
	}
	
	// 获取当前页的章节
	pageSections := allSections[start:end]
	
	return &model.ListSectionsResponse{
		Sections:    pageSections,
		Total:       total,
		CurrentPage: req.Page,
		TotalPages:  totalPages,
		HasNext:     req.Page < totalPages,
		HasPrev:     req.Page > 1,
	}, nil
}

// SelectSection 选择章节
func (s *Service) SelectSection(req *model.SelectSectionRequest) (*model.SelectSectionResponse, error) {
	ctx := context.Background()
	
	// 检查章节是否存在
	exists, err := s.sectionDAO.SectionExists(ctx, req.SectionName)
	if err != nil {
		return nil, fmt.Errorf("检查章节存在性失败: %w", err)
	}
	
	if !exists {
		return &model.SelectSectionResponse{
			IsSuccess: false,
		}, fmt.Errorf("章节 '%s' 不存在", req.SectionName)
	}
	
	// 获取章节详情
	section, err := s.sectionDAO.GetSection(ctx, req.SectionName)
	if err != nil {
		return nil, fmt.Errorf("获取章节详情失败: %w", err)
	}
	
	// 设置当前章节
	s.currentSection = req.SectionName
	
	return &model.SelectSectionResponse{
		Selected:  *section,
		WordCount: len(section.Words),
		IsSuccess: true,
	}, nil
}

// GetCurrentSection 获取当前章节
func (s *Service) GetCurrentSection() string {
	if s.currentSection == "" {
		return "未选择章节"
	}
	return s.currentSection
}

// SetCurrentSection 设置当前章节
func (s *Service) SetCurrentSection(section string) error {
	s.currentSection = section
	fmt.Printf("设置当前章节为: %s\n", section)
	return nil
}