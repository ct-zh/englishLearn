package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ct-zh/englishLearn/model"
)

// SectionDAOImpl 章节DAO实现
type SectionDAOImpl struct {
	filePath string
	mutex    sync.RWMutex
}

// NewSectionDAO 创建新的SectionDAO实例
func NewSectionDAO(dataDir string) SectionDAOInterface {
	return &SectionDAOImpl{
		filePath: filepath.Join(dataDir, "sections.json"),
	}
}

// loadData 加载JSON数据
func (s *SectionDAOImpl) loadData() (model.WordsDataDAO, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 检查文件是否存在
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// 文件不存在，返回空数据
		return make(model.WordsDataDAO), nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	var wordsData model.WordsDataDAO
	if len(data) == 0 {
		// 文件为空，返回空数据
		return make(model.WordsDataDAO), nil
	}

	if err := json.Unmarshal(data, &wordsData); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return wordsData, nil
}

// saveData 保存JSON数据
func (s *SectionDAOImpl) saveData(data model.WordsDataDAO) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 确保目录存在
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化JSON失败: %w", err)
	}

	if err := os.WriteFile(s.filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// CreateSection 创建章节
func (s *SectionDAOImpl) CreateSection(ctx context.Context, section *model.SectionEntity) error {
	data, err := s.loadData()
	if err != nil {
		return err
	}

	// 检查章节是否已存在
	if _, exists := data[section.Name]; exists {
		return fmt.Errorf("章节 '%s' 已存在", section.Name)
	}

	// 添加新章节
	data[section.Name] = section.Words
	if data[section.Name] == nil {
		data[section.Name] = make([]model.WordEntity, 0)
	}

	return s.saveData(data)
}

// GetSection 根据名称获取章节
func (s *SectionDAOImpl) GetSection(ctx context.Context, name string) (*model.SectionEntity, error) {
	data, err := s.loadData()
	if err != nil {
		return nil, err
	}

	words, exists := data[name]
	if !exists {
		return nil, fmt.Errorf("章节 '%s' 不存在", name)
	}

	return &model.SectionEntity{
		Name:  name,
		Words: words,
	}, nil
}

// UpdateSection 更新章节
func (s *SectionDAOImpl) UpdateSection(ctx context.Context, name string, section *model.SectionEntity) error {
	data, err := s.loadData()
	if err != nil {
		return err
	}

	// 检查章节是否存在
	if _, exists := data[name]; !exists {
		return fmt.Errorf("章节 '%s' 不存在", name)
	}

	// 如果需要重命名章节
	if section.Name != name {
		// 检查新名称是否已存在
		if _, exists := data[section.Name]; exists {
			return fmt.Errorf("章节 '%s' 已存在", section.Name)
		}
		// 删除旧名称，添加新名称
		delete(data, name)
		data[section.Name] = section.Words
	} else {
		// 只更新单词列表
		data[name] = section.Words
	}

	return s.saveData(data)
}

// DeleteSection 删除章节
func (s *SectionDAOImpl) DeleteSection(ctx context.Context, name string) error {
	data, err := s.loadData()
	if err != nil {
		return err
	}

	// 检查章节是否存在
	if _, exists := data[name]; !exists {
		return fmt.Errorf("章节 '%s' 不存在", name)
	}

	// 删除章节
	delete(data, name)

	return s.saveData(data)
}

// ListSections 列出所有章节
func (s *SectionDAOImpl) ListSections(ctx context.Context) ([]model.SectionEntity, error) {
	data, err := s.loadData()
	if err != nil {
		return nil, err
	}

	sections := make([]model.SectionEntity, 0, len(data))
	for name, words := range data {
		sections = append(sections, model.SectionEntity{
			Name:  name,
			Words: words,
		})
	}

	return sections, nil
}

// SectionExists 检查章节是否存在
func (s *SectionDAOImpl) SectionExists(ctx context.Context, name string) (bool, error) {
	data, err := s.loadData()
	if err != nil {
		return false, err
	}

	_, exists := data[name]
	return exists, nil
}

// AddWordToSection 向章节添加单词
func (s *SectionDAOImpl) AddWordToSection(ctx context.Context, sectionName string, word model.WordEntity) error {
	data, err := s.loadData()
	if err != nil {
		return err
	}

	// 检查章节是否存在
	words, exists := data[sectionName]
	if !exists {
		return fmt.Errorf("章节 '%s' 不存在", sectionName)
	}

	// 检查单词是否已存在
	for _, existingWord := range words {
		if existingWord.W == word.W {
			return fmt.Errorf("单词 '%s' 在章节 '%s' 中已存在", word.W, sectionName)
		}
	}

	// 添加单词
	data[sectionName] = append(words, word)

	return s.saveData(data)
}

// RemoveWordFromSection 从章节移除单词
func (s *SectionDAOImpl) RemoveWordFromSection(ctx context.Context, sectionName string, wordText string) error {
	data, err := s.loadData()
	if err != nil {
		return err
	}

	// 检查章节是否存在
	words, exists := data[sectionName]
	if !exists {
		return fmt.Errorf("章节 '%s' 不存在", sectionName)
	}

	// 查找并移除单词
	newWords := make([]model.WordEntity, 0, len(words))
	found := false
	for _, word := range words {
		if word.W != wordText {
			newWords = append(newWords, word)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("单词 '%s' 在章节 '%s' 中不存在", wordText, sectionName)
	}

	data[sectionName] = newWords

	return s.saveData(data)
}