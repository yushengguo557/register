package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yushengguo557/register/api"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())   // 使用Logger中间件
	r.Use(gin.Recovery()) // 使用Recovery中间件
	// docs.SwaggerInfo.BasePath = "http://localhost:8080/"
	// r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// http://127.0.0.1:8080/docs/index.html

	r.GET("/ping", api.Pong)

	apiv1 := r.Group("/v1") // 路由组
	{
		apiv1.POST("/wxMiniLogin", api.WXMiniProgramLogin) // 微信小程序登录
		apiv1.POST("/getPhoneNumber", api.GetPhoneNumber)  // 微信小程序获取手机号码

		apiv1.POST("/wxLogin", api.WXLogin)
	}
	return r
}
