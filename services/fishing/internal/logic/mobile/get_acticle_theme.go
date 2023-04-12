package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) GetArticleTheme(ctx context.Context, request *mobile_fishing.GetArticleThemeRequest, response *mobile_fishing.GetArticleThemeResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "GetArticleTheme",
	})
	if request.ThemeId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetArticleTheme, err := models.GetArticleTheme(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.ThemeId)
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetArticleTheme",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	if GetArticleTheme == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	response.Info = &mobile_fishing.ArticleTheme{
		ThemeId:    GetArticleTheme.Id,
		CategoryId: GetArticleTheme.CategoryId,
		TagId:      GetArticleTheme.TagId,
		Title:      GetArticleTheme.Title,
		Desc:       GetArticleTheme.Desc,
		BannerlUrl: GetArticleTheme.BannerlUrl,
		CreateTime: GetArticleTheme.CreateTime,
		ImageUrl:   GetArticleTheme.ImgUrl,
	}
	return nil
}
