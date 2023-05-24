package main

import (
	"loginTest/api"
	"loginTest/controller"
	"loginTest/middleware"

	"github.com/gin-gonic/gin"
)

// 创建路由

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.POST("/api/auth/apiTest", api.ApiTest)
	r.POST("/api/auth/post", controller.Post)
	r.POST("/api/auth/browse", controller.Browse)
	r.POST("/api/auth/updateSave", controller.UpdateSave)
	r.POST("/api/auth/updateLike", controller.UpdateLike)
	r.POST("/api/auth/deletePost", controller.DeletePost)
	r.POST("/api/auth/submitReport", controller.SubmitReport)
	r.POST("/api/auth/showDetails", controller.ShowDetails)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	return r
}
