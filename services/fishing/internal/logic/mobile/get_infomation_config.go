package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) GetInfomationConfig(ctx context.Context, request *mobile_fishing.GetInfomationConfigRequest, response *mobile_fishing.GetInfomationConfigResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "GetInfomationConfig",
	})

	ListBanners, err := models.ListBanners(func(session *xorm.Session) *xorm.Session {
		return session.Where("is_del = ?", 0).
			Where("banner_type = ?", mobile_fishing.BannerType_banner_type_infomation).
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
	response.ListInfomationCategory = append(response.ListInfomationCategory, &mobile_fishing.InfomationCategory{
		InfomationCategoryType: mobile_fishing.InfomationCategoryType_infomation_type_new,
		Name:                   "最新",
	})
	response.ListInfomationCategory = append(response.ListInfomationCategory, &mobile_fishing.InfomationCategory{
		InfomationCategoryType: mobile_fishing.InfomationCategoryType_infomation_type_nearby,
		Name:                   "附近",
	})
	response.ListInfomationCategory = append(response.ListInfomationCategory, &mobile_fishing.InfomationCategory{
		InfomationCategoryType: mobile_fishing.InfomationCategoryType_infomation_type_free,
		Name:                   "免费",
	})
	response.ListInfomationCategory = append(response.ListInfomationCategory, &mobile_fishing.InfomationCategory{
		InfomationCategoryType: mobile_fishing.InfomationCategoryType_infomation_type_charge,
		Name:                   "有偿",
	})
	response.ListInfomationCategory = append(response.ListInfomationCategory, &mobile_fishing.InfomationCategory{
		InfomationCategoryType: mobile_fishing.InfomationCategoryType_infomation_type_follow,
		Name:                   "关注",
	})
	return nil
}
