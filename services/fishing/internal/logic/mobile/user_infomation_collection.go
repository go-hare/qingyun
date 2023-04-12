package mobile

import (
	"context"
	"github.com/realmicro/realmicro/common/util/time"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"strconv"
	"xorm.io/xorm"
)

func (m *MobileFishingService) UserInfomationCollection(ctx context.Context, request *mobile_fishing.UserInfomationCollectionRequest, response *mobile_fishing.UserInfomationCollectionResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "UserInfomationCollection",
	})
	if request.UserId < 1 || request.CollectionId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	tableName := "user_collection_" + strconv.Itoa(int(request.UserId)%models.UserCollectionNum)
	GetUserCollection, err := models.GetUserCollection(func(session *xorm.Session) *xorm.Session {
		return session.Table(tableName).Where("user_collection_type = ?", mobile_fishing.UserColltionType_user_colltion_type_infomation).
			Where("collection_id = ?", request.CollectionId).Where("user_id = ?", request.UserId)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUserLike",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	GetInfomation, err := models.GetInfomation(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.CollectionId).
			Where("is_del = ?", 0).
			Where("infomation_status = ?", mobile_fishing.InfomationStatus_infomation_status_author_pass)
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
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_user_collection_error)
		return nil
	}

	if request.IfCollection { //点赞
		if GetUserCollection != nil {
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_user_collection_error)
			return nil
		}
		if err = models.CreateUserCollection(&models.UserCollection{
			UserId:             request.UserId,
			CollectionId:       request.CollectionId,
			UserCollectionType: int64(mobile_fishing.UserColltionType_user_colltion_type_infomation),
			CollectionTime:     time.Now(),
		}); err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "CreateUserCollection",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}

		models.UpdateInfomation(func(session *xorm.Session) *xorm.Session {
			return session.Where("id = ?", request.CollectionId).Incr("collect_count")
		}, GetInfomation)

	} else { //取消点赞
		if GetUserCollection == nil {
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
			return nil
		}
		if err = models.DeleteUserCollection(func(session *xorm.Session) *xorm.Session {
			return session.Table(tableName).Where("id = ?", GetUserCollection.Id)
		}); err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "DeleteUserCollection",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}

		models.UpdateInfomation(func(session *xorm.Session) *xorm.Session {
			return session.Where("id = ?", request.CollectionId).Decr("collect_count")
		}, GetInfomation)
	}

	return nil
}
