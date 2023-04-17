package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"qingyun/services/wechat/internal/common/response"
	"qingyun/services/wechat/internal/wechat"
	"strings"
)

var logger = log.WithFields(log.Fields{
	"Module": "Service",
	"Method": "GetUser",
})

// CheckAppKeyIsLoggedInMiddleware 检查AppKey是否已登录微信
func CheckAppKeyIsLoggedInMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appKey := ctx.Request.Header.Get("AppKey")
		// TODO 从数据库判断AppKey是否存在
		// 如果不是登录请求，判断AppKey是否有效
		flag := true
		if !strings.Contains(ctx.Request.RequestURI, "login") {
			if err := wechat.CheckBot(appKey); err != nil {
				logger.WithFields(log.Fields{
					"ErrorType": "Database",
					"Function":  "AppKey",
				}).Error(err)
				response.FailWithMessage("AppKey预检失败："+err.Error(), ctx)
			}
		}
		if flag {
			ctx.Next()
		} else {
			ctx.Abort()
		}
	}
}

// CheckAppKeyExistMiddleware 检查是否有appKey
func CheckAppKeyExistMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appKey := ctx.Request.Header.Get("AppKey")
		// 先判断AppKey是不是传了
		if len(appKey) < 1 {
			response.FailWithMessage("AppKey为必传参数", ctx)
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}
