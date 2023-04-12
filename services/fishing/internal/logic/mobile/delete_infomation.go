package mobile

import (
	"context"
	"github.com/realmicro/realmicro/common/util/time"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) DeleteInfomation(ctx context.Context, request *mobile_fishing.DeleteInfomationRequest, response *mobile_fishing.DeleteInfomationResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "DeleteInfomation",
	})

	if request.UserId < 1 || request.InfomationId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetInfomation, err := models.GetInfomation(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId).Where("id = ?", request.InfomationId).Where("is_del = ?", 0)
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetInfomation",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if GetInfomation == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	if err = models.UpdateInfomation(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.InfomationId).Cols("is_del,del_time")
	}, &models.Infomation{IsDel: 1, DelTime: time.Now()}); err != nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if err = models.UpdateUserExtend(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId).Decr("infomation_count")
	}, &models.UserExtend{}); err != nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	return nil
}
