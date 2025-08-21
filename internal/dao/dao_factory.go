package dao

import (
	"fmt"
	"path/filepath"
	"github.com/ct-zh/englishLearn/config"
)

// DAOFactory DAO工厂
type DAOFactory struct {
	dataFilePath string
	sectionDAO   SectionDAOInterface
	config       *config.Config // 添加配置引用
}

// NewDAOFactory 创建新的DAO工厂
func NewDAOFactory(dataFilePath string) *DAOFactory {
	return &DAOFactory{
		dataFilePath: dataFilePath,
	}
}

// NewDAOFactoryWithConfig 创建带配置的DAO工厂
func NewDAOFactoryWithConfig(cfg *config.Config) *DAOFactory {
	return &DAOFactory{
		dataFilePath: cfg.DataFilePath,
		config:       cfg,
	}
}

// GetSectionDAO 获取章节DAO实例
func (f *DAOFactory) GetSectionDAO() SectionDAOInterface {
	if f.sectionDAO == nil {
		f.sectionDAO = NewSectionDAO(filepath.Dir(f.dataFilePath))
	}
	return f.sectionDAO
}

// ProvideDAOFactory 提供DAO工厂实例 (Wire Provider)
func ProvideDAOFactory(cfg *config.Config) *DAOFactory {
	return NewDAOFactoryWithConfig(cfg)
}

// ProvideSectionDAO 提供章节DAO实例 (Wire Provider)
func ProvideSectionDAO(factory *DAOFactory) SectionDAOInterface {
	return factory.GetSectionDAO()
}

// GetDataFilePath 获取数据文件路径
func (f *DAOFactory) GetDataFilePath() string {
	return f.dataFilePath
}

// ReloadDataFile 重新加载数据文件
func (f *DAOFactory) ReloadDataFile(newFilePath string) error {
	// 如果有配置引用，更新配置
	if f.config != nil {
		if err := f.config.UpdateDataFilePath(newFilePath); err != nil {
			return err
		}
		newFilePath = f.config.DataFilePath // 使用配置处理后的路径
	}
	
	// 清理现有的DAO实例
	f.sectionDAO = nil
	
	// 更新数据文件路径
	f.dataFilePath = newFilePath
	
	// 重新初始化SectionDAO（延迟初始化）
	// 下次调用GetSectionDAO时会自动创建新实例
	
	return nil
}

// RollbackDataFile 回滚数据文件
func (f *DAOFactory) RollbackDataFile() error {
	if f.config == nil {
		return fmt.Errorf("无法回滚：缺少配置引用")
	}
	
	// 使用配置的回滚功能
	if err := f.config.RollbackDataFilePath(); err != nil {
		return err
	}
	
	// 清理现有的DAO实例
	f.sectionDAO = nil
	
	// 更新数据文件路径
	f.dataFilePath = f.config.DataFilePath
	
	return nil
}

// GetCurrentFileInfo 获取当前文件信息
func (f *DAOFactory) GetCurrentFileInfo() (map[string]interface{}, error) {
	if f.config != nil {
		return f.config.GetFileInfo()
	}
	
	// 如果没有配置引用，返回基本信息
	info := make(map[string]interface{})
	info["path"] = f.dataFilePath
	return info, nil
}