package dao

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ct-zh/englishLearn/model"
)

type WordInputMode interface {
	Prepare()
}

func ReadWordsFromFile(filename string) (model.Words, error) {
	var words model.Words
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return words, nil
		}
		return words, err
	}
	err = json.Unmarshal(file, &words)
	return words, err
}

func WriteWordsToFile(filename string, words model.Words) error {
	data, err := json.MarshalIndent(words, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}
