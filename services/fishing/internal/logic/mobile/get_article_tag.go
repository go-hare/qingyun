package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) GetArticleTag(ctx context.Context, request *mobile_fishing.GetArticleTagRequest, response *mobile_fishing.GetArticleTagResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "GetArticleTag",
	})
	if request.TagId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetArticleTag, err := models.GetArticleTag(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.TagId)
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetArticleTag",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	if GetArticleTag == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	response.Info = &mobile_fishing.ArticleTag{
		TagId:     GetArticleTag.Id,
		TagName:   GetArticleTag.Name,
		UsedCount: GetArticleTag.UsedCount,
		ViewCount: GetArticleTag.ViewCount,
	}
	return nil
}
