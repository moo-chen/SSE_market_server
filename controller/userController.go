package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"loginTest/common"
	"loginTest/dto"
	"loginTest/model"
	"loginTest/response"
	"loginTest/util"
	"net/http"
)

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("phone = ?", telephone).First(&user)
	if user.UserID != 0 {
		return true
	}
	return false
}

func isEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	if user.UserID != 0 {
		return true
	}
	return false
}

func Register(c *gin.Context) {
	// 连接数据库
	db := common.GetDB()
	// 从前端获取信息
	var requestUser = model.User{}
	c.Bind(&requestUser)
	//获取参数
	name := requestUser.Name
	telephone := requestUser.Phone
	password := requestUser.Password
	email := requestUser.Email

	//若使用postman等工具，写法如下：
	// name := c.PostForm("name")
	// telephone := c.PostForm("telephone")
	// password := c.PostForm("password")
	// PostForm中的参数与在postman中发送信息的名字要一致

	//验证数据
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "手机号必须为11位!!!")
		println(telephone)
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码不能少于6位")
		return
	}
	//如果名称没有传，给一个10位的随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	check := false
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			check = true
		}
	}
	if !check {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱不符合要求")
		return
	}
	// 判断手机号是否存在
	if isTelephoneExist(db, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该手机号已存在")
		return
	}
	if isEmailExist(db, email) {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该邮箱已存在")
		return
	}
	// 创建用户
	hasedPassword := password
	//hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//if err != nil {
	//	response.Response(c, http.StatusInternalServerError, 500, nil, "密码加密错误")
	//	return
	//}
	// 创建新用户结构体
	newUser := model.User{
		Name:     name,
		Phone:    telephone,
		Password: string(hasedPassword),
		Email:    email,
	}
	// 将结构体传进Create函数即可在数据库中添加一条记录
	// 其他的增删查改功能参见postController里的updateLike函数
	db.Create(&newUser)
	if newUser.UserID == 0 {
		response.Response(c, http.StatusInternalServerError, 400, nil, "userID为0")
		return
	}
	//发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 400, nil, "系统异常")

		log.Printf("token generate error: %v", err)
		return
	}
	//返回结果
	response.Success(c, gin.H{"token": token}, "注册成功")
}

type loginuser struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null"`
	Phone    string `gorm:"type:varchar(11);not null"`
	Password string `gorm:"size:255;not null"`
}

func Login(c *gin.Context) {
	db := common.GetDB()
	var requestUser = loginuser{}
	c.Bind(&requestUser)
	//获取参数
	telephone := requestUser.Phone
	password := requestUser.Password
	//数据验证
	if len(telephone) == 0 {
		msg := "手机号为空！" + telephone
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, msg)
		return
	}
	if len(telephone) != 11 {
		msg := "手机号必须为11位" + telephone
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, msg)
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码不能少于6位")
		return
	}
	//判断手机号是否存在
	var user model.User
	db.Where("phone = ?", telephone).First(&user)
	if user.UserID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "用户不存在")
		return
	}
	//判断密码是否正确
	//if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	if user.Password != password {
		response.Response(c, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 400, nil, "系统异常")

		log.Printf("token generate error: %v", err)
		return
	}
	//返回结果
	response.Success(c, gin.H{"token": token}, "登录成功")
}

type modifyUser struct {
	Phone    string `gorm:"type:varchar(11);not null"`
	Password string `gorm:"size:255;not null"`
}

func ModifyPassword(c *gin.Context) {
	fmt.Println("Successfully deliver!")
	db := common.GetDB()
	var user model.User
	var inputUser modifyUser
	c.Bind(&inputUser)
	phone := inputUser.Phone
	password := inputUser.Password
	db.Where("phone = ?", phone).First(&user)
	fmt.Println(phone)
	fmt.Println(password)
	if user.UserID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "用户不存在")
		return
	}
	user.Password = password
	db.Save(&user)
	response.Success(c, gin.H{"data": user}, "修改密码成功")
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}
