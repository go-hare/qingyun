package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) CancelUserFollow(ctx context.Context, request *mobile_fishing.CancelUserFollowRequest, response *mobile_fishing.CancelUserFollowResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "CancelUserFollow",
	})

	if request.UserId < 1 || request.FollowUserId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetUserFollow, err := models.GetUserFollow(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId).
			Where("follow_user_id = ?", request.FollowUserId)
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUserFollow",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if GetUserFollow == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	if err = models.UpdateUserExtend(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId).Decr("follow_count")
	}, &models.UserExtend{}); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "UpdateUserExtend",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if err = models.UpdateUserExtend(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.FollowUserId).Decr("fans_count")
	}, &models.UserExtend{}); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "UpdateUserExtend",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if err = models.DeleteUserFollow(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", GetUserFollow.Id)
	}); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "DeleteUserFollow",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	return nil
}
