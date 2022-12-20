package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DataSource struct {
		Type  string `json:",optional,default=mysql"`
		Table string `json:",optional,default=version"`
		Url   string
	}

	Path string `json:",default=./command"`
}

func Load(configFile string) *Config {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Errorf("error:%+v", err)
		return nil
	}
	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Errorf("error:%+v", err)
		return nil
	}
	return &conf
}
