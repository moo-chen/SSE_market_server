package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"log"
	"loginTest/common"
	"loginTest/config"
	"loginTest/route"
	"net/http"
	"os/exec"
	"time"
)

func Copy() {
	// 数据库连接信息
	dbHost := viper.GetString("datasource.host")
	dbPort := viper.GetInt("datasource.port")
	dbUser := viper.GetString("datasource.username")
	dbPassword := viper.GetString("datasource.password")
	dbName := viper.GetString("datasource.database")

	// 备份目录
	backupDir := "/Users/michael/Documents/backup"

	c := cron.New()
	c.AddFunc("@daily", func() {
		backupFile := fmt.Sprintf("%s/backup_%s.sql", backupDir, time.Now().Format("2006-01-02 15:04:05"))
		cmd := exec.Command("mysqldump", fmt.Sprintf("-h%s", dbHost), fmt.Sprintf("-P%d", dbPort), fmt.Sprintf("-u%s", dbUser), fmt.Sprintf("-p%s", dbPassword), dbName, "--result-file="+backupFile)
		err := cmd.Run()
		if err != nil {
			log.Println("备份失败:", err)
			return
		}
		log.Println("备份成功:", backupFile)
	})
	c.Start()
}

var r *gin.Engine

func main() {
	config.InitConfig()
	Copy()
	db := common.InitDB()
	rds := common.RedisInit()
	defer rds.Close()
	defer db.Close()
	r = gin.Default()

	// 使用 http.FileServer 文件服务器处理 "/uploads/" 开头的请求，
	// 文件服务器获取文件的位置在 "./public" 文件夹下。
	r.StaticFS("/uploads", http.Dir("./public/uploads"))

	route.CollectRoute(r)
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
