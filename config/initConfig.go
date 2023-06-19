package config

import (
	"github.com/spf13/viper"
	"os"
	"strings"
)

// 使用viper从配置文件中读取配置
func InitConfig() {
	wordDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath("config")
	wordDir = strings.TrimRight(wordDir, "/test")
	viper.AddConfigPath(wordDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
