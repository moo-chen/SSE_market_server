package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"loginTest/common"
	"loginTest/dto"
	"loginTest/model"
	"loginTest/response"
	"net/http"
)

type User struct {
	UserID    int    `json:"-"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Name      string `json:"name"`
	Num       int    `json:"num"`
	Profile   string `json:"-"`
	Intro     string `json:"-"`
	IDpass    bool   `json:"IDpass"`
	Ban       bool   `json:"ban"`
	Punishnum int    `json:"punishnum"`
}

type Username struct {
	Name string
}

type modifyUser struct {
	Account   string
	Password1 string
	Password2 string
}

type Check struct {
	Name   string
	Phone  string
	IdPass int
}

// 登录
func AdminLogin(ctx *gin.Context) {
	db := common.DB
	var requesrAdmin model.Admin
	ctx.Bind(&requesrAdmin)

	account := requesrAdmin.Account
	password := requesrAdmin.Password

	var admin model.Admin
	db.Where("account = ?", account).First(&admin)
	if admin.AdminID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "管理员账号不存在")
		return
	}
	if password != admin.Password {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "密码错误")
		return
	}
	//发放token
	token, err := common.ReleaseToken_admin(admin)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 400, nil, "系统异常")

		log.Printf("token generate error: %v", err)
		return
	}
	//返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func AdminInfo(c *gin.Context) {
	admin, _ := c.Get("admin")
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"admin": dto.ToAdminDto(admin.(model.Admin))}})
}

// 输出所有用户
func ShowFilterUsers(ctx *gin.Context) {
	db := common.DB
	var userList []User
	var requestInfo = Check{}
	ctx.Bind(&requestInfo)

	name := requestInfo.Name
	phone := requestInfo.Phone
	idPass := requestInfo.IdPass

	if phone != "" {
		db = db.Model(&model.User{}).Where("phone = ?", phone)
	}
	if name != "" {
		db = db.Model(&model.User{}).Where("name like ?", name+"%")
	}
	if idPass != -1 {
		db = db.Model(&model.User{}).Where("idPass = ?", idPass)
	}

	db.Find(&userList)
	response.Success(ctx, gin.H{"data": userList}, "Successfully show all users")
}

// 更改是否审查
func PassUsers(ctx *gin.Context) {
	fmt.Println("start to pass")
	db := common.DB
	var username = Username{}
	ctx.Bind(&username)
	name := username.Name
	fmt.Println(username)
	fmt.Println(name)
	db.Model(&model.User{}).Where("name = ?", name).Update("IDpass", true)
	var user model.User
	db.Where("name = ?", name).Find(&user)
	response.Success(ctx, gin.H{"data": user}, "Successfully pass user")
}

// 添加管理员
func AddAdmin(ctx *gin.Context) {
	fmt.Println("Start to add")
	db := common.DB
	var newAdmin modifyUser
	ctx.Bind(&newAdmin)
	account := newAdmin.Account
	pass1 := newAdmin.Password1
	pass2 := newAdmin.Password2
	var admin model.Admin
	if account == "" {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "账号不能为空")
		return
	}

	if pass1 == "" || pass2 == "" {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "密码不能为空")
		return
	}

	if pass1 != pass2 {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "两次密码不同，请重新输入")
		return
	}

	if len(pass1) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "密码必须大于等于6位")
		return
	}

	db.Where("account = ?", account).First(&admin)
	fmt.Println(admin.Account)
	if admin.Account != "" {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "管理员账号已存在")
		return
	}

	addAdmin := model.Admin{
		Account:  account,
		Password: pass1,
	}
	db.Create(&addAdmin)

	response.Success(ctx, gin.H{"data": addAdmin}, "添加管理员成功")
}

// 修改密码
func ChangeAdminPassword(ctx *gin.Context) {
	db := common.DB
	var admin modifyUser
	var newAdmin model.Admin

	ctx.Bind(&admin)
	account := admin.Account
	pass1 := admin.Password1
	pass2 := admin.Password2
	if pass1 != pass2 {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "两次密码不同，请重新输入")
		return
	}

	db.Where("Account = ?", account).First(&newAdmin)
	if newAdmin.Account == "" {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "管理员不存在")
		return
	}
	addAdmin := model.Admin{
		Account:  account,
		Password: pass1,
	}

	db.Save(&addAdmin)
	response.Success(ctx, gin.H{"data": newAdmin}, "成功修改管理员密码")
}

// 注销用户账号
func DeleteUser(ctx *gin.Context) {
	db := common.DB
	var user = Username{}
	ctx.Bind(&user)

	fmt.Println(user)
	name := user.Name
	fmt.Println(name)
	var checkUser model.User
	db.Where("name = ?", name).First(&checkUser)
	if checkUser.UserID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "未找到该用户")
		return
	}

	db.Delete(&checkUser)
	response.Response(ctx, http.StatusOK, 200, nil, "成功删除该用户")
}

// 注销管理员
func DeleteAdmin(ctx *gin.Context) {
	db := common.DB
	var user model.Admin
	ctx.Bind(&user)

	account := user.Account
	var checkUser model.Admin
	db.Where("account = ?", account).First(&checkUser)
	if checkUser.Account == "" {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "未找到该管理员")
		return
	}

	db.Delete(&checkUser)
	response.Success(ctx, nil, "成功删除该管理员")
}
