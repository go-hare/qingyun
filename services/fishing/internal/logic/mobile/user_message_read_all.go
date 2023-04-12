package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"strconv"
	"xorm.io/xorm"
)

func (m *MobileFishingService) UserMessageReadAll(ctx context.Context, request *mobile_fishing.UserMessageReadAllRequest, response *mobile_fishing.UserMessageReadAllResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "UserMessageReadAll",
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
	tableName := "message_" + strconv.Itoa(int(request.UserId)%models.MessageNum)
	if err = models.UpdateMessage(func(session *xorm.Session) *xorm.Session {
		return session.Table(tableName).Where("user_id = ?", request.UserId).Where("read_status = ?", 0)
	}, &models.Message{ReadStatus: 1}); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "UpdateMessage",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if err = models.UpdateUserExtend(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.UserId)
	}, &models.UserExtend{UnReadCount: 0}); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "UpdateUser",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	return nil
}
