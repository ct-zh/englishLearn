package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// 通用config配置文件读取方法

type Config struct {
	Dsn string `yaml:"dsn"`
} // 通用配置

func New() *Config {
	cfg := &Config{}
	cfgFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(cfgFile, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
