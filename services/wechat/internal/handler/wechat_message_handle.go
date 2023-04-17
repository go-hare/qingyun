package handler

import (
	"github.com/eatmoreapple/openwechat"
	"qingyun/services/wechat/plugins"
)

func HandleMessage(bot *openwechat.Bot) {
	// 定义一个处理器
	dispatcher := openwechat.NewMessageMatchDispatcher()
	// 设置为异步处理
	dispatcher.SetAsync(true)
	// 处理消息为已读
	dispatcher.RegisterHandler(checkIsCanRead, setTheMessageAsRead)

	//// 默认启用插件
	plugins.ChangePluginStatus(true)
	// 注册插件
	dispatcher.RegisterHandler(plugins.WeChatPluginInstance.CheckIsOpen,
		plugins.WeChatPluginInstance.OpenGPT)

	//// 注册文本消息处理函数
	dispatcher.OnText(textMessageHandle)
	//// 注册图片消息处理器
	//dispatcher.OnImage(imageMessageHandle)
	//// 注册表情包消息处理器
	//dispatcher.OnEmoticon(emoticonMessageHandle)
	//// APP消息处理
	//dispatcher.OnMedia(appMessageHandle)
	//// 保存消息
	//dispatcher.RegisterHandler(checkNeedSave, saveToDb)
	//// 未定义消息处理
	//dispatcher.RegisterHandler(checkIsOther, otherMessageHandle)

	// 注册消息处理器
	bot.MessageHandler = dispatcher.AsMessageHandler()
}
