package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DataSource struct {
		Type  string `mapstructure:"Type,optional,default=mysql"`
		Table string `mapstructure:"Table,optional,default=version"`
		Url   string
	}

	Path string `json:",default=./sqlCommandVersion"`
}

func MustLoad(configFile string) *Config {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.SetDefault("DataSource.Type", "mysql")
	viper.SetDefault("DataSource.Table", "version")
	viper.SetDefault("Path", "./sqlCommandVersion")

	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}

	return &conf
}
