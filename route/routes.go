package route

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
	noauth.POST("/api/auth/uploadphotos", controller.UploadPhotos)
	noauth.POST("/api/auth/uploadavatar", controller.UploadAvatar)
	auth.POST("/api/auth/getavatar", controller.GetAvatar)
	auth.POST("/api/auth/getInfo", controller.GetInfo)
	auth.POST("/api/auth/updateUserInfo", controller.UpdateUserInfo)
	auth.POST("/api/auth/uploadZip", controller.UploadZip)
	auth.POST("/api/auth/submitFeedback", controller.SubmitFeedback)
	auth.GET("/api/auth/getAllFeedback", controller.GetAllFeedback)
	auth.GET("/api/auth/getNotice", controller.GetNotice)
	auth.PATCH("api/auth/readNotice/:noticeID", controller.ReadNotice)
	noauth.POST("/api/auth/modifyPassword", controller.ModifyPassword)
	noauth.POST("/api/auth/validateEmail", controller.ValidateEmail)
	noauth.POST("/api/auth/identityValidate", controller.IdentityValidate)
	r.POST("/api/auth/passUsers", controller.PassUsers)
	r.POST("/api/auth/addAdmin", controller.AddAdmin)
	r.POST("/api/auth/changePassword", controller.ChangeAdminPassword)
	r.POST("/api/auth/deleteUser", controller.DeleteUser)
	r.POST("/api/auth/deleteAdmin", controller.DeleteAdmin)
	r.POST("/api/auth/showUsers", controller.ShowFilterUsers)
	r.POST("/api/auth/adminLogin", controller.AdminLogin)
	r.GET("/api/auth/admininfo", middleware.AuthMiddleware_admin(), controller.AdminInfo)
	return r
}
