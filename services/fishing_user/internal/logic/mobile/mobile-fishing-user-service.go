package mobile

import (
	"qingyun/services/fishing_user/internal/svc"
)

type MobileFishingUserService struct {
	SvcContext *svc.ServiceContext
}

func NewMobileFishingUserService(ctx *svc.ServiceContext) (mobileUserService *MobileFishingUserService) {
	mobileUserService = &MobileFishingUserService{
		SvcContext: ctx,
	}
	return
}
