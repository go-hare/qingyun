package mobile

import (
	"context"
	"github.com/realmicro/realmicro/common/util/time"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) DeleteActicle(ctx context.Context, request *mobile_fishing.DeleteArticleRequest, response *mobile_fishing.DeleteArticleResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "DeleteActicle",
	})

	if request.UserId < 1 || request.ArticleId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetArticle, err := models.GetArticle(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId).
			Where("id = ?", request.ArticleId).
			Where("is_del = ?", 0)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetArticle",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if GetArticle == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	if err = models.UpdateArticle(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.ArticleId).Cols("is_del,del_time")
	}, &models.Article{IsDel: 1, DelTime: time.Now()}); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "UpdateArticle",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	return nil
}
