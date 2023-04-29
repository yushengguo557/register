package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yushengguo557/register/server"
)

// WXMiniProgramLogin 微信小程序登录
func WXMiniProgramLogin(c *gin.Context) {
	// 1.解析 小程序 wx.login() 后发送到服务器的 code
	code := c.PostForm("code")

	// 2.实例化服务 并 登录
	s := server.Server{}
	ret, err := s.WXMiniProgramLogin(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ret)
}

func GetPhoneNumber(c *gin.Context) {
	// 1.从请求头中获取 token
	ACCESS_TOKEN := c.GetHeader("Authorization")
	if ACCESS_TOKEN == "" {
		// 如果请求头中没有 token，则返回错误信息
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
		return
	}

	// 2.解析 小程序 wx.login() 后发送到服务器的 code
	code := c.PostForm("code")

	// 3.发送请求获取手机号码
	s := server.Server{}
	if err := s.GetPhoneNumber(code, ACCESS_TOKEN); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// WXLogin 微信登录
func WXLogin(c *gin.Context) {}
