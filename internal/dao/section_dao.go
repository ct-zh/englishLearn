package dao

import (
	"context"
	"github.com/ct-zh/englishLearn/model"
)

// SectionDAOInterface 章节DAO接口
type SectionDAOInterface interface {
	// CreateSection 创建章节
	CreateSection(ctx context.Context, section *model.SectionEntity) error
	
	// GetSection 根据名称获取章节
	GetSection(ctx context.Context, name string) (*model.SectionEntity, error)
	
	// UpdateSection 更新章节
	UpdateSection(ctx context.Context, name string, section *model.SectionEntity) error
	
	// DeleteSection 删除章节
	DeleteSection(ctx context.Context, name string) error
	
	// ListSections 列出所有章节
	ListSections(ctx context.Context) ([]model.SectionEntity, error)
	
	// SectionExists 检查章节是否存在
	SectionExists(ctx context.Context, name string) (bool, error)
	
	// AddWordToSection 向章节添加单词
	AddWordToSection(ctx context.Context, sectionName string, word model.WordEntity) error
	
	// RemoveWordFromSection 从章节移除单词
	RemoveWordFromSection(ctx context.Context, sectionName string, wordText string) error
}