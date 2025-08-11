package utils

import (
	"encoding/json"
	"os"
)

// FileExists 检查文件是否存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// ReadJSONFile 读取JSON文件
func ReadJSONFile(filename string, v interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// WriteJSONFile 写入JSON文件
func WriteJSONFile(filename string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// CreateDirIfNotExist 如果目录不存在则创建
func CreateDirIfNotExist(dir string) error {
	if !FileExists(dir) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}