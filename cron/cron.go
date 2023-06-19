package main

import (
	"fmt"
	"log"
	"math"
	"time"

	// "loginTest/common"
	"loginTest/model"

	// "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
)

type MyJob struct {
	Heat float64
	ctx  *cron.Cron
}

func (j MyJob) Run() {
	RefreshHeat(j.ctx)
	fmt.Println(j.Heat, time.Now())
}

func main() {
	// 执行启动时的操作
	fmt.Println("项目启动中...")
	// 创建 cron 实例并添加任务
	c := cron.New()
	c.Start()
	// 在这里改每日需要刷新热榜的时间
	c.AddFunc("36 17 * * *", func() { MyJob{Heat: 0}.Run() })
	select {}
}

func RefreshHeat(c *cron.Cron) {
	// 初始化数据库连接
	db, err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/ssemarket")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// 获取所有帖子的热度值
	var heatValues []float64
	db.Model(&model.Post{}).Pluck("heat", &heatValues)

	// 计算每个帖子的新热度（开根号原来的值），
	for i := range heatValues {
		heatValue := math.Sqrt(heatValues[i])
		// 更新热度值
		db.Model(&model.Post{}).Where("heat = ?", heatValues[i]).Update("heat", heatValue)
	}

	// 在Terminal输出，表示更新成功
	fmt.Println("Refreshing heat values at:", time.Now())
}
