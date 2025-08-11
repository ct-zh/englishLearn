package sections

import (
	"fmt"
	"github.com/ct-zh/englishLearn/model"
)

// Service sections业务逻辑服务
type Service struct {
	// TODO: 后续可以注入DAO层依赖
}

// NewService 创建sections服务实例
func NewService() *Service {
	return &Service{}
}

// AddWord 添加单词
func (s *Service) AddWord(req *model.AddWordRequest) error {
	// TODO: 实现添加单词的业务逻辑
	fmt.Printf("添加单词: %s (%s) - \n", req.Word, req.Translation)
	fmt.Printf("章节: %s\n", req.Section)
	
	// 这里应该调用DAO层保存数据
	// 暂时只是打印信息
	return nil
}

// ListWords 获取单词列表
func (s *Service) ListWords(req *model.ListWordsRequest) ([]model.Word, error) {
	// TODO: 实现获取单词列表的业务逻辑
	fmt.Printf("获取章节 %s 的单词列表\n", req.Section)
	fmt.Printf("分页: 第%d页，每页%d条\n", req.Page, req.Size)
	
	// 这里应该调用DAO层查询数据
	// 暂时返回空列表
	return []model.Word{}, nil
}

// RandomWords 随机练习单词
func (s *Service) RandomWords(req *model.RandomWordsRequest) ([]model.Word, error) {
	// TODO: 实现随机练习的业务逻辑
	fmt.Printf("从章节 %s 随机选择 %d 个单词进行练习\n", req.Section, req.Count)
	
	// 这里应该调用DAO层随机查询数据
	// 暂时返回空列表
	return []model.Word{}, nil
}

// SearchWord 搜索单词
func (s *Service) SearchWord(req *model.SearchWordRequest) ([]model.Word, error) {
	// TODO: 实现搜索单词的业务逻辑
	fmt.Printf("在章节 %s 中搜索关键词: %s\n", req.Section, req.Keyword)
	
	// 这里应该调用DAO层搜索数据
	// 暂时返回空列表
	return []model.Word{}, nil
}

// GetCurrentSection 获取当前章节
func (s *Service) GetCurrentSection() string {
	// TODO: 从配置或会话中获取当前章节
	return "day1"
}

// SetCurrentSection 设置当前章节
func (s *Service) SetCurrentSection(section string) error {
	// TODO: 保存当前章节到配置或会话
	fmt.Printf("设置当前章节为: %s\n", section)
	return nil
}