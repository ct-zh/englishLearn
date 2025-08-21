package dao

import (
	"path/filepath"
	"github.com/ct-zh/englishLearn/config"
)

// DAOFactory DAO工厂
type DAOFactory struct {
	dataDir    string
	sectionDAO SectionDAOInterface
}

// NewDAOFactory 创建新的DAO工厂
func NewDAOFactory(dataDir string) *DAOFactory {
	return &DAOFactory{
		dataDir: dataDir,
	}
}

// GetSectionDAO 获取章节DAO实例
func (f *DAOFactory) GetSectionDAO() SectionDAOInterface {
	if f.sectionDAO == nil {
		f.sectionDAO = NewSectionDAO(f.dataDir)
	}
	return f.sectionDAO
}

// ProvideDAOFactory 提供DAO工厂实例 (Wire Provider)
func ProvideDAOFactory(cfg *config.Config) *DAOFactory {
	return NewDAOFactory(filepath.Dir(cfg.DataFilePath))
}

// ProvideSectionDAO 提供章节DAO实例 (Wire Provider)
func ProvideSectionDAO(factory *DAOFactory) SectionDAOInterface {
	return factory.GetSectionDAO()
}

// GetDataFilePath 获取数据文件路径
func (f *DAOFactory) GetDataFilePath() string {
	return filepath.Join(f.dataDir, "sections.json")
}