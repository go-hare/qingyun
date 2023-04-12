package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) GetUser(ctx context.Context, request *mobile_fishing.GetUserRequest, response *mobile_fishing.GetUserResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "GetUser",
	})

	if request.UserId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
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
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if getUser == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	if getUser.UserStatus == int64(mobile_fishing.UserStatus_user_status_block) {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_block_error)
		return nil
	}
	response.Info = &mobile_fishing.User{
		UserId:          getUser.Id,
		Avatar:          getUser.AvatarUrl,
		Mobile:          getUser.Mobile,
		NickName:        getUser.NickName,
		UserType:        mobile_fishing.UserType(getUser.UserType),
		CoinAmount:      getUser.CoinAmount,
		PointAmount:     getUser.PointAmount,
		Motto:           getUser.Motto,
		CoinPayTotal:    getUser.CoinPayTotal,
		CoinTotal:       getUser.CoinTotal,
		CoinRefillTotal: getUser.CoinRefillTotal,
	}
	return nil
}
