package mobile

import (
	"qingyun/services/fishing/internal/svc"
)

type MobileFishingService struct {
	SvcContext *svc.ServiceContext
}

func NewMobileFishingService(ctx *svc.ServiceContext) (mobileService *MobileFishingService) {
	mobileService = &MobileFishingService{
		SvcContext: ctx,
	}
	return
}
