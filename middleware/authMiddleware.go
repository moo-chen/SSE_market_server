package middleware

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
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
		var user model.User
		err = db.QueryRow("SELECT * FROM User WHERE ID = ?",userId).Scan(&user.ID, &user.Name, &user.Telephone, &user.Password)
		//用户不存在
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		} else if err != nil {
			log.Fatal(err)
		}
		//用户存在，将user的信息写入上下文
		ctx.Set("user",user)
		ctx.Next()
	}
}
