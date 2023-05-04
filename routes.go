package main

import (
	"loginTest/controller"
	"loginTest/middleware"
	"loginTest/api"
	"github.com/gin-gonic/gin"
)

// 创建路由

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.POST("/api/auth/apiTest",api.ApiTest)
	r.POST("/api/auth/post",controller.Post)
	r.POST("/api/auth/browse",controller.Browse)
	r.POST("/api/auth/updateLike",controller.UpdateLike)
	r.GET("/api/auth/info", middleware.AuthMiddleware(),controller.Info)
	return r
}
