package plugins

import "github.com/eatmoreapple/openwechat"

const (
	isOpne = true
	proxy  = "http://127.0.0.1:7890"
	gptKey = "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"
)

type weChatPlugin struct {
	isOpen bool
}

var WeChatPluginInstance *weChatPlugin

func init() {
	if WeChatPluginInstance == nil {
		WeChatPluginInstance = &weChatPlugin{false}
	}
}

// ChangePluginStatus 修改插件状态
func ChangePluginStatus(isOpen bool) {
	WeChatPluginInstance = &weChatPlugin{isOpen}
}

// CheckIsOpen 检查是否开启微信插件
func (p weChatPlugin) CheckIsOpen(message *openwechat.Message) bool {
	return p.isOpen
}
