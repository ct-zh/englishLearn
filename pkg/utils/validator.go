package utils

import (
	"errors"
	"strings"
)

// ValidateWord 验证单词数据
func ValidateWord(word, chinese, phrase string) error {
	if strings.TrimSpace(word) == "" {
		return errors.New("单词不能为空")
	}
	if strings.TrimSpace(chinese) == "" {
		return errors.New("中文释义不能为空")
	}
	// phrase 可以为空，所以不验证
	return nil
}

// ValidateSection 验证章节名称
func ValidateSection(section string) error {
	if strings.TrimSpace(section) == "" {
		return errors.New("章节名称不能为空")
	}
	return nil
}

// IsValidMenuOption 验证菜单选项
func IsValidMenuOption(option string, maxOption int) bool {
	if len(option) != 1 {
		return false
	}
	
	if option[0] >= '1' && option[0] <= '0'+byte(maxOption) {
		return true
	}
	
	// 支持 'q' 或 '0' 退出
	return option == "q" || option == "Q" || option == "0"
}