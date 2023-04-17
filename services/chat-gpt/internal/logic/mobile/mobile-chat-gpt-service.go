package mobile

import (
	"qingyun/services/chat-gpt/internal/svc"
)

type MobileChatGptService struct {
	SvcContext *svc.ServiceContext
}

func NewMobileChatGptService(ctx *svc.ServiceContext) (mobileService *MobileChatGptService) {
	mobileService = &MobileChatGptService{
		SvcContext: ctx,
	}
	return
}
