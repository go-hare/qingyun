package mobile

import (
	"qingyun/services/book-keeping/internal/svc"
)

type MobileBookKeepingService struct {
	SvcContext *svc.ServiceContext
}

func NewMobileBookKeepingService(ctx *svc.ServiceContext) (mobileBookKeepingService *MobileBookKeepingService) {
	mobileBookKeepingService = &MobileBookKeepingService{
		SvcContext: ctx,
	}
	return
}
