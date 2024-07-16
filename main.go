package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/ct-zh/englishLearn/config"
	"github.com/ct-zh/englishLearn/dao"
	"github.com/ct-zh/englishLearn/model"
)

const (
	loadWords = 1
)

type Program struct {
	Dao *dao.Dao
}

func main() {
	cfg := config.New()
	p := Program{
		Dao: dao.Init(cfg),
	}

	for {
		fmt.Println("请输入数字1进入单词输入模式:")
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case loadWords:
			p.enterWordInputMode()
		default:
			fmt.Println("无效输入，请重新输入。")
		}
	}
}

// 进入单词输入模式
func (p *Program) enterWordInputMode() {
	var words model.Words
	const filename = "words.json"

	words, err := dao.ReadWordsFromFile(filename)
	if err != nil {
		fmt.Println("读取文件出错:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("请输入一个单词:")
	inputWord, _ := reader.ReadString('\n')
	inputWord = inputWord[:len(inputWord)-1] // 去掉换行符

	exists := false
	for i, word := range words.Words {
		if word.Word == inputWord {
			exists = true
			fmt.Println("单词已存在，请输入新的短语:")
			inputPhrase, _ := reader.ReadString('\n')
			inputPhrase = inputPhrase[:len(inputPhrase)-1] // 去掉换行符
			words.Words[i].Phrase = inputPhrase
			words.Words[i].UpdatedAt = time.Now().Format(time.RFC3339)
			break
		}
	}

	if !exists {
		fmt.Println("请输入该单词的短语:")
		inputPhrase, _ := reader.ReadString('\n')
		inputPhrase = inputPhrase[:len(inputPhrase)-1] // 去掉换行符
		newWord := model.Word{
			Word:      inputWord,
			Phrase:    inputPhrase,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}
		words.Words = append(words.Words, newWord)
	}

	err = dao.WriteWordsToFile(filename, words)
	if err != nil {
		fmt.Println("写入文件出错:", err)
	} else {
		fmt.Println("单词已保存。")
	}
}
