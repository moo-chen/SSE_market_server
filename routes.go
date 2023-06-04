package main

import (
	"github.com/gin-gonic/gin"
	"loginTest/api"
	"loginTest/controller"
	"loginTest/middleware"
)

// 创建路由

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.POST("/api/auth/modifyPassword", controller.ModifyPassword)
	r.POST("/api/auth/apiTest", api.ApiTest)
	r.POST("/api/auth/post", controller.Post)
	r.POST("/api/auth/browse", controller.Browse)
	r.POST("/api/auth/updateLike", controller.UpdateLike)
	r.POST("/api/auth/validateEmail", controller.ValidateEmail)
	r.POST("/api/auth/showDetails", controller.ShowDetails)
	r.POST("/api/auth/identityValidate", controller.IdentityValidate)
	r.POST("/api/auth/passUsers", controller.PassUsers)
	r.POST("/api/auth/addAdmin", controller.AddAdmin)
	r.POST("/api/auth/modifyAdminPassword", controller.ChangeAdminPassword)
	r.POST("/api/auth/deleteUser", controller.DeleteUser)
	r.POST("/api/auth/deleteAdmin", controller.DeleteAdmin)
	r.POST("/api/auth/showUsers", controller.ShowFilterUsers)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	return r
}
