package common

import (
	"fmt"
	"loginTest/model"
	_ "github.com/alexbrainman/odbc"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	// 使用viper从配置文件中读取数据库配置
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
			username,
			password,
			host,
			port,
			database,
			charset)
	db, err:= gorm.Open(driverName,args)
	if err != nil{
		panic("failed to connect database, err: "+err.Error())
	}
	// 若没有相应数据库，运行时将根据对应结构体自动创建数据库
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Post{})
	db.AutoMigrate(&model.Like{})
	DB = db
	return db
}

func GetDB() *gorm.DB{
	return DB
}