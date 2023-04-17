package logic

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"github.com/realmicro/realmicro/common/util/time"
	log "github.com/sirupsen/logrus"
	"qingyun/services/wechat/internal/common/proto"
	"qingyun/services/wechat/internal/common/response"
	"qingyun/services/wechat/internal/wechat"
	"qingyun/services/wechat/models"
	"unicode/utf8"
	"xorm.io/xorm"
)

// loginUrlResponse
// @description: 获取登录URL返回结构体
type loginUrlResponse struct {
	Uuid string `json:"uuid"`
	Url  string `json:"url"`
}

// GetLoginUrlHandle 获取登录扫码连接
func GetLoginUrlHandle(ctx *gin.Context) {
	appKey := ctx.Request.Header.Get("AppKey")
	// 获取一个微信机器人对象
	bot := wechat.InitWechatBotHandle()
	// 已扫码回调
	bot.ScanCallBack = func(body openwechat.CheckLoginResponse) {
		log.Infof("[%v]已扫码", appKey)
	}

	// 设置登录成功回调
	bot.LoginCallBack = func(body openwechat.CheckLoginResponse) {
		log.Infof("[%v]登录成功", appKey)
	}

	// 获取登录二维码链接
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	uuid, err := bot.Caller.GetLoginUUID()
	if err != nil {
		log.Errorf("获取登录二维码失败: %v", err.Error())
		response.FailWithMessage("获取登录二维码失败："+err.Error(), ctx)
		return
	}
	log.Infof("获取到uuid: %v", uuid)
	// 拼接URL
	url := fmt.Sprintf("https://login.weixin.qq.com/qrcode/%s", uuid)

	// 保存Bot到实例
	wechat.SetBot(appKey, bot)
	// 返回数据
	response.OkWithData(loginUrlResponse{Uuid: uuid, Url: url}, ctx)
}

// LoginHandle 登录
func LoginHandle(ctx *gin.Context) {
	appKey := ctx.Request.Header.Get("AppKey")
	uuid := ctx.Query("uuid")
	if utf8.RuneCountInString(uuid) < 1 {
		response.FailWithMessage("uuid为必传参数", ctx)
		return
	}
	//usePush := ctx.Query("usePush") // 是否使用免扫码登录
	//isPush := usePush == "1" || usePush == "true" || usePush == "yes"
	GetWechatBot, err := models.GetWechatBot(func(session *xorm.Session) *xorm.Session {
		return session.Where("app_key = ?", appKey)
	})

	if err != nil {
		log.Errorf("登录失败: %v", err)
		return
	}

	if GetWechatBot == nil {
		models.CreateWechatBot(&models.WechatBot{
			AppKey:     appKey,
			CreateTime: time.Now(),
		})
	}

	// 获取Bot对象
	bot := wechat.GetBot(appKey)
	if bot == nil {
		response.FailWithMessage("请先获取登录二维码", ctx)
		return
	}

	// 定义登录数据缓存
	storage := proto.NewMysqlHotReloadStorage(appKey)

	// 热登录
	var opts []openwechat.BotLoginOption
	opts = append(opts, openwechat.NewRetryLoginOption()) // 热登录失败使用扫码登录，适配第一次登录的时候无热登录数据
	//opts = append(opts, openwechat.NewSyncReloadDataLoginOption(10*time.Minute)) // 十分钟同步一次热登录数据

	// 登录
	if err := bot.HotLogin(storage, opts...); err != nil {
		log.Errorf("登录失败: %v", err)
		response.FailWithMessage("登录失败："+err.Error(), ctx)
		return
	}

	// 获取登录用户信息
	user, err := bot.GetCurrentUser()
	if err != nil {
		log.Errorf("获取登录用户信息失败: %v", err.Error())
		response.FailWithMessage("获取登录用户信息失败："+err.Error(), ctx)
		return
	}
	models.UpdateWechatBot(func(session *xorm.Session) *xorm.Session {
		return session.Where("app_key = ?", appKey)
	}, &models.WechatBot{
		UserName:  user.UserName,
		NickName:  user.NickName,
		LoginTime: time.Now(),
		IfLogin:   1,
	})
	log.Infof("当前登录用户：%v", user.NickName)
	response.OkWithMessage("登录成功", ctx)
}
