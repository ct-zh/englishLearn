package config

import (
	"path/filepath"
)

// Config 应用配置结构体
type Config struct {
	DataFilePath string // JSON数据文件路径
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		DataFilePath: filepath.Join("data", "sections.json"),
	}
}

// LoadConfig 加载配置
// TODO: 实现从配置文件加载配置的功能
func LoadConfig() (*Config, error) {
	// 目前返回默认配置
	return DefaultConfig(), nil
}