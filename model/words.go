package model

type Word struct {
	Word      string `json:"word"`
	Phrase    string `json:"phrase"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Words struct {
	Words []Word `json:"words"`
}
