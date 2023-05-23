package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"loginTest/api"
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

type registerUser struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null"`
	Phone    string `gorm:"type:varchar(11);not null"`
	Password string `gorm:"size:255;not null"`
	Email    string `gorm:"type:varchar(11);not null"`
	ValiCode string `gorm:"type:varchar(10);not null"`
}

func Register(c *gin.Context) {
	// 连接数据库
	db := common.GetDB()
	// 从前端获取信息
	var requestUser = registerUser{}
	c.Bind(&requestUser)
	//获取参数
	name := requestUser.Name
	telephone := requestUser.Phone
	password := requestUser.Password
	email := requestUser.Email
	valiCode := requestUser.ValiCode

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
	rds := common.MyRedis
	ctx := context.Background()
	correctValiCode, err := rds.Get(ctx, email).Result()
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
		return
	}
	fmt.Println(correctValiCode)
	fmt.Println(valiCode)
	if correctValiCode != valiCode {
		response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
		return
	}
	// 创建用户
	// hasedPassword := password
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "密码加密错误")
		return
	}
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
	rds.Del(ctx, email)
}

type identity struct {
	Email    string `gorm:"type:varchar(11);not null"`
	ValiCode string `gorm:"type:varchar(10);not null"`
}

func IdentityValidate(c *gin.Context) {
	db := common.GetDB()
	// 从前端获取信息
	var requestUser = identity{}
	c.Bind(&requestUser)
	//获取参数
	email := requestUser.Email
	valiCode := requestUser.ValiCode

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
	if !isEmailExist(db, email) {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该邮箱不存在")
		return
	}
	rds := common.MyRedis
	ctx := context.Background()
	correctValiCode, err := rds.Get(ctx, email).Result()
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
		return
	}
	fmt.Println(correctValiCode)
	fmt.Println(valiCode)
	if correctValiCode != valiCode {
		response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
		return
	}

	var newUser model.User
	db.Where("email = ?", email).First(&newUser)
	//发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 400, nil, "系统异常")

		log.Printf("token generate error: %v", err)
		return
	}
	//返回结果
	response.Success(c, gin.H{"token": token}, "身份验证成功")
	rds.Del(ctx, email)
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
	Phone     string `gorm:"type:varchar(11);not null"`
	Password  string `gorm:"size:255;not null"`
	Password2 string `gorm:"size:255;not null"`
}

func ModifyPassword(c *gin.Context) {
	fmt.Println("Successfully deliver!")
	db := common.GetDB()
	var user model.User
	var inputUser modifyUser
	c.Bind(&inputUser)
	phone := inputUser.Phone
	password := inputUser.Password
	password2 := inputUser.Password2
	fmt.Println(phone)
	fmt.Println(password)
	fmt.Println(password2)
	if password != password2 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码输入不一致")
		return
	}
	if !isTelephoneExist(db, phone) {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "未找到电话")
		return
	}
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

type requestEmail struct {
	Email string `gorm:"type:varchar(11);not null"`
	Mode  int
}

func ValidateEmail(c *gin.Context) {
	var request requestEmail
	c.Bind(&request)
	email := request.Email
	mode := request.Mode
	db := common.DB
	if mode == 0 && isEmailExist(db, email) {
		response.Response(c, http.StatusBadRequest, 400, gin.H{"data": email, "mode": mode}, "邮箱已注册")
		return
	}
	fmt.Println("email is ", email)
	if email == "" {
		response.Response(c, http.StatusBadRequest, 400, gin.H{"data": email}, "邮箱获取错误")
		return
	}
	err := api.SendEmail(email)
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, gin.H{"data": email}, "发送邮箱错误")
		return
	}
	response.Success(c, gin.H{"data": email}, "邮箱发送成功")
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}
