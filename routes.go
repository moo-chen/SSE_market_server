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
	r.Use(middleware.AuthMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.POST("/api/auth/apiTest", api.ApiTest)
	r.POST("/api/auth/post", controller.Post)
	r.POST("/api/auth/browse", controller.Browse)
	r.POST("/api/auth/updateLike", controller.UpdateLike)
	r.POST("/api/auth/showDetails", controller.ShowDetails)
	r.POST("api/auth/showPcomments", controller.GetComments)
	r.POST("api/auth/postPcomment", controller.PostPcomment)
	r.POST("api/auth/postCcomment", controller.PostCcomment)
	r.POST("api/auth/updateCcommentLike", controller.UpdateCcommentLike)
	r.POST("api/auth/updatePcommentLike", controller.UpdatePcommentLike)

	return r

}
