package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"loginTest/common"
	"loginTest/dto"
	"loginTest/model"
	"loginTest/response"
	"loginTest/util"
	"net/http"
)

func isTelephoneExist(db *sql.DB, telephone string) bool {
	has, err := db.Query("SELECT * FROM User WHERE Telephone = ?", telephone)
	if err != nil {
		log.Fatal(err)
	}
	if has.Next() {
		return true
	} else {
		return false
	}
}
func Register(c *gin.Context) {
	db := common.GetDB()
	var requestUser = model.User{}
	c.Bind(&requestUser)
	//获取参数
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	//验证数据
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	//如果名称没有传，给一个10位的随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	// 判断手机号是否存在
	if isTelephoneExist(db, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "该手机号已存在")
		return
	}
	// 创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "密码加密错误")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	stmt, err := db.Prepare("INSERT INTO User(Name, Telephone, Password) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(newUser.Name, newUser.Telephone, newUser.Password)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(newUser.Name, newUser.Telephone, newUser.Password)
	response.Success(c, nil, "注册成功")
}
func Login(c *gin.Context) {
	db := common.GetDB()
	var requestUser = model.User{}
	c.Bind(&requestUser)
	//获取参数
	telephone := requestUser.Telephone
	password := requestUser.Password
	//数据验证
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	//判断手机号是否存在
	var user model.User
	err := db.QueryRow("SELECT * FROM User WHERE Telephone = ?", telephone).Scan(&user.ID, &user.Name, &user.Telephone, &user.Password)
	if err == sql.ErrNoRows {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	} else if err != nil {
		log.Fatal(err)
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(c, http.StatusBadRequest, 422, nil, "密码错误")
		return
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 422, nil, "系统异常")

		log.Printf("token generate error: %v", err)
		return
	}
	//返回结果
	response.Success(c, gin.H{"token": token}, "登录成功")
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}