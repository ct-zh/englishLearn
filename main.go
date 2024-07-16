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

	ChoiceInterface map[int]interface{}
}

func main() {
	cfg := config.New()
	p := Program{
		Dao: dao.Init(cfg),
		ChoiceInterface: map[int]interface{}{
			loadWords: dao.NewWordsFromJson(cfg),
		},
	}

	for {
		fmt.Println("请输入数字1进入单词输入模式:")
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case loadWords:
			p.enterWordInputMode(p.ChoiceInterface[loadWords].(dao.WordInputMode))
		default:
			fmt.Println("无效输入，请重新输入。")
		}
	}
}

// 进入单词输入模式
func (p *Program) enterWordInputMode(wIMode dao.WordInputMode) {
	wIMode.Prepare()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("请输入一个单词:")
	inputWord, _ := reader.ReadString('\n')
	inputWord = inputWord[:len(inputWord)-1] // 去掉换行符

	exists := false
	word, err := wIMode.FindWord(inputWord)
	if err != nil {
		if err != dao.NotFoundErr {
			panic(fmt.Errorf("查找单词出错: %v", err))
		}
	}
	if word != nil {
		exists = true
	}

	if exists {
		fmt.Println("单词已存在，请输入新的短语:")
		inputPhrase, _ := reader.ReadString('\n')
		inputPhrase = inputPhrase[:len(inputPhrase)-1] // 去掉换行符
		word.Phrase = inputPhrase
		word.UpdatedAt = time.Now().Format(time.DateTime)
		if err := wIMode.Update(word); err != nil {
			panic(fmt.Errorf("更新单词出错: %v", err))
		}
	}

	if !exists {
		fmt.Println("请输入该单词的短语:")
		inputPhrase, _ := reader.ReadString('\n')
		inputPhrase = inputPhrase[:len(inputPhrase)-1] // 去掉换行符
		newWord := &model.Word{
			Word:      inputWord,
			Phrase:    inputPhrase,
			CreatedAt: time.Now().Format(time.DateTime),
			UpdatedAt: time.Now().Format(time.DateTime),
		}
		if err := wIMode.Insert(newWord); err != nil {
			panic(fmt.Errorf("写入单词出错: %v", err))
		}
	}

	fmt.Println("单词已保存。")
}
