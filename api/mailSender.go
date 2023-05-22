package api

import (
	"context"
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"loginTest/common"
	"loginTest/someConst"
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

func formVcode(ctx string) (string, string) {
	rand.Seed(time.Now().Unix()) // unix 时间戳，秒
	data := rand.Int() % 1000000
	vcode := ""
	for i := 1; i <= 6; i++ {
		vcode += strconv.Itoa(data % 10)
		data = data / 10
	}
	ctx = strings.Replace(ctx, "vcode", vcode, -1)
	fmt.Println(ctx)
	fmt.Println(vcode)
	return ctx, vcode
}

func saveVcode(vcode, receiver string) {
	rds := common.MyRedis
	ctx := context.Background()
	_, err := rds.Get(ctx, receiver).Result()
	if err == nil {
		rds.Del(ctx, receiver)
	}
	rds.Set(ctx, receiver, vcode, 0)
}

func SendEmail(receiver string) error {
	e := email.NewEmail()
	senderString := someConst.ValidateSender
	senderString = strings.Replace(senderString, "emailUsername", someConst.EmailUsername, -1)
	e.From = senderString

	e.To = []string{receiver}
	e.Subject = someConst.ValidateTitle
	text := someConst.ValidateText

	text, vcode := formVcode(text)
	saveVcode(vcode, receiver)
	e.Text = []byte(text)

	err := e.Send(someConst.Addr, smtp.PlainAuth("", someConst.EmailUsername, someConst.Password, someConst.Host))
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Send Successfully")
	return nil
}
