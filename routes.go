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
	// 这里把路由分成了两组，其中auth是需要token验证的，也就是需要用户登录，noauth是不需要token的，也就是不需要用户登录。
	auth := r.Group("")
	noauth := r.Group("")

	auth.Use(middleware.AuthMiddleware())
	noauth.POST("/api/auth/register", controller.Register)
	noauth.POST("/api/auth/login", controller.Login)
	auth.POST("/api/auth/apiTest", api.ApiTest)
	auth.POST("/api/auth/post", controller.Post)
	auth.POST("/api/auth/browse", controller.Browse)
	auth.POST("/api/auth/updateLike", controller.UpdateLike)
	auth.POST("/api/auth/showDetails", controller.ShowDetails)
	auth.POST("api/auth/showPcomments", controller.GetComments)
	auth.POST("api/auth/postPcomment", controller.PostPcomment)
	auth.POST("api/auth/postCcomment", controller.PostCcomment)
	auth.POST("api/auth/updateCcommentLike", controller.UpdateCcommentLike)
	auth.POST("api/auth/updatePcommentLike", controller.UpdatePcommentLike)
	auth.POST("/api/auth/updateSave", controller.UpdateSave)
	auth.POST("/api/auth/deletePost", controller.DeletePost)
	auth.POST("/api/auth/submitReport", controller.SubmitReport)
	auth.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	return r

}
