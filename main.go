package main

import (
	// _ "github.com/alexbrainman/odbc"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"loginTest/common"
	"net/http"
	"os"
)

func main() {
	InitConfig()
	db := common.InitDB()
	rds := common.RedisInit()
	defer rds.Close()
	defer db.Close()
	r := gin.Default()

	// 使用 http.FileServer 文件服务器处理 "/uploads/" 开头的请求，
	// 文件服务器获取文件的位置在 "./public" 文件夹下。
	r.StaticFS("/uploads", http.Dir("./public/uploads"))

	CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
}

// 使用viper从配置文件中读取配置
func InitConfig() {
	wordDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(wordDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
