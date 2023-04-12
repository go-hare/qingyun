package mobile

import (
	"context"
	"github.com/realmicro/realmicro/common/util/time"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) UserFollow(ctx context.Context, request *mobile_fishing.UserFollowRequest, response *mobile_fishing.UserFollowResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "ListUserFollows",
	})

	if request.UserId < 1 || request.FollowUserId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetUserFollow, err := models.GetUserFollow(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId).Where("follow_user_id = ?", request.FollowUserId)
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
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_repeat_follow_error)
		return nil
	}

	user, err := models.GetUser(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.UserId).
			Where("user_status = ?", mobile_fishing.UserStatus_user_status_default).
			Select("id")
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUser",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
	}

	if user == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	followUser, err := models.GetUser(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.FollowUserId).
			Where("user_status = ?", mobile_fishing.UserStatus_user_status_default).
			Select("id")
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUser",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
	}

	if followUser == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	if err = models.CreateUserFollow(&models.UserFollow{
		UserId:       user.Id,
		FollowUserId: followUser.Id,
		FollowTime:   time.Now(),
	}); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "CreateUserFollow",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	models.UpdateUserExtend(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", user.Id).Incr("follow_count")
	}, &models.UserExtend{})
	models.UpdateUserExtend(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", followUser.Id).Incr("fans_count")
	}, &models.UserExtend{})
	return nil
}
