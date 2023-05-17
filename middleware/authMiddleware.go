// 返回token中间件，不用管，要看也行。模板写法来着，至于为什么，我也不是很懂
package middleware

import (
	"github.com/gin-gonic/gin"
	"loginTest/common"
	"loginTest/model"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")
		//若token为空或token格式不正确
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}
		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}
		//验证通过后获取claims中的userId
		userId := claims.UserId
		db := common.GetDB()
		user := model.User{}
		db.Where("userID =?", userId).First(&user)
		//用户不存在
		if user.Userid == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}
		//用户存在，将user的信息写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}
