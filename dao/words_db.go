package dao

import (
	"github.com/ct-zh/englishLearn/model"
)

type WordsFromDB struct {
	Dao *Dao
}

func NewWordsFromDB(dao *Dao) *WordsFromDB {
	return &WordsFromDB{Dao: dao}
}

func (w *WordsFromDB) Prepare() {

}

func (w *WordsFromDB) FindWord(word string) (*model.Word, error) {
	m := &model.Word{}
	err := w.Dao.Db.Where("word = ?", word).Find(m).Error
	return m, err
}

func (w *WordsFromDB) Update(word *model.Word) error {
	return w.Dao.Db.Save(word).Error
}

func (w *WordsFromDB) Insert(word *model.Word) error {
	return w.Dao.Db.Create(word).Error
}

var _ WordInputMode = (*WordsFromDB)(nil)
