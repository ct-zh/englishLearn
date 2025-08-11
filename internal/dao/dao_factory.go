package dao

import (
	"path/filepath"
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

// GetSectionDAO 获取SectionDAO实例
func (f *DAOFactory) GetSectionDAO() SectionDAOInterface {
	if f.sectionDAO == nil {
		f.sectionDAO = NewSectionDAO(f.dataDir)
	}
	return f.sectionDAO
}

// GetDataFilePath 获取数据文件路径
func (f *DAOFactory) GetDataFilePath() string {
	return filepath.Join(f.dataDir, "sections.json")
}