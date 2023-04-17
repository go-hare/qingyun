package logic

import "github.com/gin-gonic/gin"

// initLoginRoute 初始化登录路由信息
func InitLoginRoute(app *gin.Engine) {
	// 获取登录二维码
	app.GET("/login", GetLoginUrlHandle)
	// 检查登录状态
	app.POST("/login", LoginHandle)
}
