package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) GetUserMessageUnCount(ctx context.Context, request *mobile_fishing.GetUserMessageUnCountRequest, response *mobile_fishing.GetUserMessageUnCountResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "GetUserMessageUnCount",
	})

	if request.UserId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	getUser, err := models.GetUserExtend(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId)
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUser",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	if getUser == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	//
	//if getUser.UserStatus == int64(mobile_fishing.UserStatus_user_status_block) {
	//	mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_block_error)
	//	return nil
	//}
	response.UnCount = getUser.UnReadCount
	return nil
}
