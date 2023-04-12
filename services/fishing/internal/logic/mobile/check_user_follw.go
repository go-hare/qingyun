package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) CheckUserFollw(ctx context.Context, request *mobile_fishing.CheckUserFollwRequest, response *mobile_fishing.CheckUserFollwResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "CheckUserFollw",
	})

	if request.UserId < 1 || request.FollowUserId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetUserFollow, err := models.GetUserFollow(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId).Where("follow_user_id = ?", request.FollowUserId).Select("id")
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUserFollow",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	if GetUserFollow != nil {
		response.IfFollow = true
	}
	return nil
}
