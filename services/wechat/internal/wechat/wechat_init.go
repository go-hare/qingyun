package wechat

import (
	"errors"
	"github.com/eatmoreapple/openwechat"
	log "github.com/sirupsen/logrus"
	"qingyun/services/wechat/internal/common/proto"
	"qingyun/services/wechat/internal/handler"
	"qingyun/services/wechat/models"
	"xorm.io/xorm"
)

var (
	// 登录用户的Bot对象
	wechatBots map[string]*openwechat.Bot
)

// InitWechatBotsMap 初始化WechatBots
func InitWechatBotsMap() {
	wechatBots = make(map[string]*openwechat.Bot)
}

// GetBot 获取Bot对象
func GetBot(appKey string) *openwechat.Bot {
	return wechatBots[appKey]
}

// SetBot 保存Bot对象
func SetBot(appKey string, bot *openwechat.Bot) {
	wechatBots[appKey] = bot
}

// CheckBot 预检AppKey是否存在登录记录且登录状态是否正常
func CheckBot(appKey string) error {
	// 判断指定AppKey是不是有登录信息
	bot := GetBot(appKey)
	if nil == bot {
		return errors.New("未获取到登录记录")
	}
	// 判断在线状态是否正常
	if !bot.Alive() {
		return errors.New("微信在线状态异常，请重新登录")
	}
	return nil
}

// InitBotWithStart 系统启动的时候从Redis加载登录信息自动登录
func InitBotWithStart() {
	keys, err := models.ListWechatBots(func(session *xorm.Session) *xorm.Session {
		return session.Where("if_login = ?", 1)
	})
	if err != nil {
		log.Error("获取Key失败")
		return
	}
	log.Infof("获取到登录用户信息数量：%v", len(keys))
	for _, key := range keys {
		// 提取出AppKey
		// 调用热登录
		log.Debugf("当前热登录用户: %v", key.AppKey)
		bot := InitWechatBotHandle()
		storage := proto.NewMysqlHotReloadStorage(key.AppKey)
		if err = bot.HotLogin(storage, openwechat.NewRetryLoginOption()); err != nil {
			log.Infof("[%v] 热登录失败，错误信息：%v", key.AppKey, err.Error())
			// 登录失败，删除热登录数据
			if err = models.UpdateWechatBot(func(session *xorm.Session) *xorm.Session {
				return session.Where("app_key = ?", key.AppKey).Cols("is_login")
			}, &models.WechatBot{IfLogin: 0}); err != nil {
				log.Errorf("[%v] Redis缓存删除失败，错误信息：%v", key, err.Error())
			}
			continue
		}
		loginUser, _ := bot.GetCurrentUser()
		log.Infof("[%v]初始化自动登录成功，用户名：%v", key.AppKey, loginUser.NickName)
		// 登录成功，写入到WechatBots
		SetBot(key.AppKey, bot)
	}
}

// InitWechatBotHandle 初始化微信机器人
func InitWechatBotHandle() *openwechat.Bot {
	bot := openwechat.DefaultBot(openwechat.Desktop)
	// 设置心跳回调
	bot.SyncCheckCallback = func(resp openwechat.SyncCheckResponse) {
		if resp.RetCode == "1100" {
			log.Errorf("微信已退出")
			// do something
		}
		switch resp.Selector {
		case "0":
			log.Debugf("正常")
		case "2", "6":
			log.Debugf("有新消息")
		case "7":
			log.Debugf("进入/离开聊天界面")
			err := bot.WebInit()
			if err != nil {
				// 短信通知一下
				// do something
				//log.Panicf("重新初始化失败: %v", err)
			}
		default:
			log.Debugf("RetCode: %s  Selector: %s", resp.RetCode, resp.Selector)
		}
	}

	// 注册消息处理函数
	handler.HandleMessage(bot)
	// 获取消息发生错误
	//bot.MessageOnError()
	// 返回机器人对象
	return bot
}
