package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) GetUserExtend(ctx context.Context, request *mobile_fishing.GetUserExtendRequest, response *mobile_fishing.GetUserExtendResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "GetUserExtend",
	})

	if request.UserId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetUserExtend, err := models.GetUserExtend(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUserExtend",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	if GetUserExtend == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	response.Info = &mobile_fishing.UserExtend{
		FollowCount:     GetUserExtend.FollowCount,
		FansCount:       GetUserExtend.FansCount,
		UnReadCount:     GetUserExtend.UnReadCount,
		InfomationCount: GetUserExtend.InfomationCount,
		ArticleCount:    GetUserExtend.ArticleCount,
		InviteCount:     GetUserExtend.InviteCount,
	}
	return nil
}
