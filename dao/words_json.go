package dao

import (
	"fmt"
	"github.com/ct-zh/englishLearn/config"
	"github.com/ct-zh/englishLearn/model"
)

type WordsFromJson struct {
	Cfg        *config.Config
	WordsCache *model.Words
}

func NewWordsFromJson(cfg *config.Config) *WordsFromJson {
	return &WordsFromJson{Cfg: cfg, WordsCache: nil}
}

var _ WordInputMode = (*WordsFromJson)(nil)

func (w *WordsFromJson) Prepare() {
	var err error
	if w.WordsCache == nil {
		w.WordsCache, err = ReadWordsFromFile(w.Cfg.JsonPath)
		if err != nil {
			panic(fmt.Errorf("读取文件出错: %v", err))
		}
	}
}

func (w *WordsFromJson) FindWord(inputWord string) (*model.Word, error) {
	var m *model.Word
	for idx, item := range w.WordsCache.Words {
		if item.Word == inputWord {
			m = &model.Word{
				Id:        int64(idx),
				Word:      item.Word,
				Phrase:    item.Phrase,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
			return m, nil
		}
	}
	return m, NotFoundErr
}

func (w *WordsFromJson) Update(word *model.Word) error {
	w.WordsCache.Words[word.Id] = word
	return WriteWordsToFile(w.Cfg.JsonPath, w.WordsCache)
}

func (w *WordsFromJson) Insert(word *model.Word) error {
	word.Id = int64(len(w.WordsCache.Words))
	w.WordsCache.Words = append(w.WordsCache.Words, word)
	return WriteWordsToFile(w.Cfg.JsonPath, w.WordsCache)
}
