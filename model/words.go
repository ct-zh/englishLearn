package model

type Word struct {
	Id        int64  `json:"id" gorm:"column:id"`
	Word      string `json:"word" gorm:"column:word"`
	Phrase    string `json:"phrase" gorm:"column:phrase"`
	CreatedAt string `json:"created_at" gorm:"column:created_at"`
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at"`
}

func (w Word) TableName() string {
	return "words"
}

type Words struct {
	Words []*Word `json:"words"`
}
