package dao

import (
	"fmt"
	"github.com/ct-zh/englishLearn/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var NotFoundErr = fmt.Errorf("not found")

type Dao struct {
	Db *gorm.DB
}

func Init(cfg *config.Config) *Dao {
	dao := &Dao{}

	if cfg.Dsn != "" {
		db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("init dao err = %v", err))
		}
		dao.Db = db
	} else {
		fmt.Printf("数据库未配置")
	}

	return dao
}
