package model

// Word 单词结构体
type Word struct {
	W      string `json:"W"`      // 原始单词
	C      string `json:"C"`      // 中文释义
	Phrase string `json:"Phrase"` // 对应短语
}

// WordsData JSON文件的数据结构
type WordsData map[string][]Word

// Section 章节信息结构体
type Section struct {
	Name  string // 章节名称
	Words []Word // 章节中的单词
}

// Response 通用响应结构体
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
