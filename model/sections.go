package model

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
	Words []Word `json:"words"`
	Total int    `json:"total"`
}