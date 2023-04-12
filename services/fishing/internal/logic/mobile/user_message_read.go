package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"strconv"
	"xorm.io/xorm"
)

func (m *MobileFishingService) UserMessageRead(ctx context.Context, request *mobile_fishing.UserMessageReadRequest, response *mobile_fishing.UserMessageReadResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "UserMessageReadAll",
	})
	if request.UserId < 1 || request.MessageId < 1 {
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
	GetMessage, err := models.GetMessage(func(session *xorm.Session) *xorm.Session {
		return session.Table(tableName).Where("id =  ?", request.MessageId).Where("user_id = ?", request.UserId).Where("read_status = ?", 0)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetMessage",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if GetMessage == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	if err = models.UpdateMessage(func(session *xorm.Session) *xorm.Session {
		return session.Table(tableName).Where("id = ?", GetMessage.Id)
	}, &models.Message{ReadStatus: 1}); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "UpdateMessage",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if err = models.UpdateUserExtend(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.UserId).Decr("un_read_count")
	}, &models.UserExtend{}); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "UpdateUser",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	return nil
}
