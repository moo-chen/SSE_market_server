package main
import (
	"loginTest/common"
	"os"
	_ "github.com/alexbrainman/odbc"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	
)
func main() {
	InitConfig()
	db := common.InitDB()
	defer db.Close()
	r := gin.Default()
	CollectRoute(r)
	port := viper.GetString("server.port")
	if port !="" {
		panic(r.Run(":"+port))
	}
	panic(r.Run())
}

// 使用viper从配置文件中读取配置
func InitConfig() {
	wordDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(wordDir+"/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}