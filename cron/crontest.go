// package main

// import (
// 	"fmt"
// 	"time"

// 	"github.com/robfig/cron/v3"
// )

// type MyJob struct {
// 	Name string
// }

// func (j MyJob) Run() {
// 	fmt.Println(j.Name, time.Now())
// }

// func main() {
// 	// 执行启动时的操作
// 	fmt.Println("项目启动中...")
// 	// 创建 cron 实例并添加任务
// 	c := cron.New(cron.WithSeconds())
// 	c.AddFunc("* * * * * *", func() { MyJob{Name: "basic"}.Run() })
// 	c.Start()

// 	select {}
// }

package main

import "fmt"

func test() {
	fmt.Println("Hello, world!")
}
