package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/yushengguo557/register/common"
)

// JWTAuthMiddleware 个中间件，用于验证 JWT token
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			// 如果请求头中没有 token，则返回错误信息
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			c.Abort()
			return
		}

		// 解析 token
		claims := &common.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return "", nil
		})
		if err != nil || !token.Valid {
			// 如果解析失败或者 token 不合法，则返回错误信息
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			c.Abort()
			return
		}

		// 将用户openid存储到上下文中，方便后续的处理
		c.Set("openid", claims.Openid)
		c.Next()
	}
}
