package model

// ===== CLI层请求/响应结构体 =====

// AddWordRequest 添加单词请求
type AddWordRequest struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
	Section     string `json:"section"`
}

// ListWordsRequest 列出单词请求
type ListWordsRequest struct {
	Section string `json:"section"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}

// RandomWordsRequest 随机单词请求
type RandomWordsRequest struct {
	Section string `json:"section"`
	Count   int    `json:"count"`
}

// SearchWordRequest 搜索单词请求
type SearchWordRequest struct {
	Keyword string `json:"keyword"`
	Section string `json:"section,omitempty"` // 可选，指定在某个章节中搜索
}

// SearchWordResponse 搜索单词响应
type SearchWordResponse struct {
	Words []WordEntity `json:"words"`
	Total int          `json:"total"`
}

// CreateSectionRequest 创建章节请求
type CreateSectionRequest struct {
	Name string `json:"name"`
}

// UpdateSectionRequest 更新章节请求
type UpdateSectionRequest struct {
	Name     string       `json:"name"`
	NewName  string       `json:"new_name,omitempty"`
	Words    []WordEntity `json:"words,omitempty"`
}

// ListSectionsResponse 列出章节响应
type ListSectionsResponse struct {
	Sections []SectionEntity `json:"sections"`
	Total    int             `json:"total"`
}

// ===== 实体结构体 =====

// WordEntity 单词实体
type WordEntity struct {
	W      string `json:"W"`      // 原始单词
	C      string `json:"C"`      // 中文释义
	Phrase string `json:"Phrase"` // 对应短语
}

// SectionEntity 章节实体
type SectionEntity struct {
	Name  string       `json:"name"`  // 章节名称
	Words []WordEntity `json:"words"` // 章节中的单词
}

// ===== DAO层数据结构体 =====

// WordsDataDAO JSON文件的数据结构
type WordsDataDAO map[string][]WordEntity

// SectionDAO 章节DAO结构体
type SectionDAO struct {
	Name  string       `json:"name"`
	Words []WordEntity `json:"words"`
}

// ===== 兼容性别名 (向后兼容) =====

// Word 单词结构体 (已废弃，使用WordEntity)
type Word = WordEntity

// Section 章节信息结构体 (已废弃，使用SectionEntity)
type Section = SectionEntity

// WordsData JSON文件的数据结构 (已废弃，使用WordsDataDAO)
type WordsData = WordsDataDAO
