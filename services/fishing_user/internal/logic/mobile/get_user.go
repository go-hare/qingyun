package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing_user/models"
	mobile_fishing_user "qingyun/services/fishing_user/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingUserService) GetUser(ctx context.Context, request *mobile_fishing_user.GetUserRequest, response *mobile_fishing_user.GetUserResponse) (err error) {
	response.Status = &mobile_fishing_user.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "GetUser",
	})

	if request.UserId < 1 {
		mobile_fishing_user.ServiceStatus(response.Status, mobile_fishing_user.StatusCode_status_param_error)
		return nil
	}

	getUser, err := models.GetUser(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.UserId)
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUser",
		}).Error(err)
		mobile_fishing_user.ServiceStatus(response.Status, mobile_fishing_user.StatusCode_status_internal_error)
	}

	if getUser == nil {
		mobile_fishing_user.ServiceStatus(response.Status, mobile_fishing_user.StatusCode_status_param_error)
		return nil
	}

	response.Info = &mobile_fishing_user.User{
		UserId:   getUser.Id,
		Avatar:   getUser.Avatar,
		Mobile:   getUser.Mobile,
		NickName: getUser.NickName,
	}
	return nil
}
