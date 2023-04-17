package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"qingyun/services/wechat/pkg/openai"
	"strings"
)

// OpenGPT
// @description: 开放式GPT-3聊天机器人
// @receiver weChatPlugin
// @param ctx
func (weChatPlugin) OpenGPT(ctx *openwechat.MessageContext) {
	// 判断是否开启了GPT-3聊天机器人
	if !isOpne {
		return
	}

	msg := ctx.Content
	// 判断是不是引用消息
	//var beforeMsg string
	if strings.HasPrefix(ctx.Content, "「") && strings.Contains(ctx.Content, "\n- - - - - - - - - - - - - - -\n") {
		// 提取出前文
		//beforeMsg = strings.Split(ctx.Content, "\n- - - - - - - - - - - - - - -\n")[0]
		msg = strings.Split(ctx.Content, "\n- - - - - - - - - - - - - - -\n")[1]
	}

	// 取出消息第一行以及剩下的内容
	msgArray := strings.Split(msg, "\n")
	if len(msgArray) < 2 || strings.ToLower(msgArray[0]) != "@openai" {
		return
	}
	// 获取提问的内容
	question := strings.Join(msgArray[1:], "\n")

	log.Debugf("ChatGPT提问内容: %s", question)
	fmt.Println("ChatGPT提问内容: %s", question)
	// 调用GPT-3聊天机器人
	// 如果配置了代理，就设置一下
	hc := http.Client{}
	if proxy != "" {
		hc.Transport = &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(proxy)
			},
		}
	}
	// 创建OpenAI客户端
	conf := openai.DefaultConfig(gptKey)
	conf.HTTPClient = &hc
	client := openai.NewClientWithConfig(conf)
	// 组装消息 TODO 懒得搞上下文联动，有想法的可以自己实现，只需要组装一下下面这个Message字段就行了，把之前的记录带过去
	request := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question,
			},
		},
	}
	// 调用聊天机器人
	resp, err := client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		_, _ = ctx.ReplyText(err.Error())
		return
	}
	log.Debugf("ChatGPT回答内容: %s", resp.Choices[0].Message.Content)
	fmt.Println("ChatGPT回答内容: %s", resp.Choices[0].Message.Content)
	// 发送回复
	if ctx.IsSendByGroup() {
		// 取出消息在群里面的发送者
		js, _ := json.Marshal(ctx)
		fmt.Println(string(js), ctx.ToUserName)
	}
	//bot := ctx.Bot()
	//ctx
	c, err := ctx.ReplyText(resp.Choices[0].Message.Content)
	fmt.Println(c, err)
}
