package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) GetActicleConfig(ctx context.Context, request *mobile_fishing.GetActicleConfigRequest, response *mobile_fishing.GetActicleConfigResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "GetActicleConfig",
	})

	if request.CategoryId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	ListBanners, err := models.ListBanners(func(session *xorm.Session) *xorm.Session {
		return session.Where("is_del = ?", 0).
			Where("banner_type = ?", mobile_fishing.BannerType_banner_type_active).
			OrderBy("weight desc")
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListBanners",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	for i := 0; i < len(ListBanners); i++ {
		response.ListBanners = append(response.ListBanners, &mobile_fishing.Banner{
			ImageUrl:   ListBanners[i].ImageUrl,
			ViewType:   mobile_fishing.BannerViewType(ListBanners[i].ViewType),
			ExternalId: ListBanners[i].ExternalId,
			Title:      ListBanners[i].Title,
			LinkUrl:    ListBanners[i].LinkUrl,
			ThemeId:    ListBanners[i].ThemeId,
		})
	}

	ListArticleThemes, err := models.ListArticleThemes(func(session *xorm.Session) *xorm.Session {
		session = session.Where("is_del = ?", 0)
		session = session.Where("category_id = ?", request.CategoryId)
		return session.OrderBy("view_count desc").Limit(6, 0)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListArticleThemes",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	tagIds := []int64{}
	for i := 0; i < len(ListArticleThemes); i++ {
		tagIds = append(tagIds, ListArticleThemes[i].TagId)
	}
	ListArticleTags, err := models.ListArticleTags(func(session *xorm.Session) *xorm.Session {
		return session.Where("is_del = ?", 0).In("id", tagIds)
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListArticleTags",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	ListArticleTagsMap := make(map[int64]*models.ArticleTag)
	for i := 0; i < len(ListArticleTags); i++ {
		ListArticleTagsMap[ListArticleTags[i].Id] = ListArticleTags[i]
	}

	for i := 0; i < len(ListArticleThemes); i++ {
		if v, ok := ListArticleTagsMap[ListArticleThemes[i].TagId]; ok {
			response.ListThemes = append(response.ListThemes, &mobile_fishing.ArticleTheme{
				ThemeId:    ListArticleThemes[i].Id,
				CategoryId: ListArticleThemes[i].CategoryId,
				TagId:      ListArticleThemes[i].TagId,
				Title:      ListArticleThemes[i].Title,
				Desc:       ListArticleThemes[i].Desc,
				//BannerlUrl: ListArticleThemes[i].BannerlUrl,
				ImageUrl:   ListArticleThemes[i].ImgUrl,
				CreateTime: ListArticleThemes[i].CreateTime,
				UsedCount:  v.UsedCount,
				ViewCount:  v.ViewCount,
			})
		}

	}
	return nil
}
