package mobile

import (
	"context"
	"github.com/realmicro/realmicro/common/util/time"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) CreateArticle(ctx context.Context, request *mobile_fishing.CreateArticleRequest, response *mobile_fishing.CreateArticleResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "CreateArticle",
	})

	if request.UserId < 1 || request.Info == nil || request.Info.Title == "" || request.Info.CategoryId < 1 || request.Info.Content == "" || request.Info.ImageUrl == "" {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	GetUser, err := models.GetUser(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.UserId).
			Where("user_status = ?", mobile_fishing.UserStatus_user_status_default)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUser",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if GetUser == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	if mobile_fishing.UserStatus(GetUser.UserStatus) == mobile_fishing.UserStatus_user_status_block {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_block_error)
		return nil
	}

	article := &models.Article{
		CategoryId:        request.Info.CategoryId,
		UserId:            request.UserId,
		Title:             request.Info.Title,
		CreateTime:        time.Now(),
		ListTags:          request.Info.TagNames,
		ImageUrl:          request.Info.ImageUrl,
		ImageWidth:        request.Info.ImageWidth,
		ImageHeight:       request.Info.ImageHeight,
		Content:           request.Info.Content,
		ListArticleTagIds: request.Info.TagIds,
	}
	if err = models.CreateArticle(article); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "Article",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	return nil
}
