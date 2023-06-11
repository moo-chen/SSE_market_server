package main

import (
	_ "github.com/alexbrainman/odbc"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
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
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServeTLS("ssl/cert.crt", "ssl/cert.key"); err != nil {
			log.Fatal("ListenAndServeTLS: ", err)
		}
	}()

	log.Printf("Server started on port 8080")
	select {}
	//port := viper.GetString("server.port")
	//if port != "" {
	//	panic(r.Run(":" + port))
	//}
	//panic(r.Run())
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
