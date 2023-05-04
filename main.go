package main

import (
	_ "github.com/alexbrainman/odbc"
	"github.com/gin-gonic/gin"
	"loginTest/common"
	
)
func main() {
	db := common.InitDB()
	defer db.Close()
	r := gin.Default()
	CollectRoute(r)
	panic(r.Run())
}
