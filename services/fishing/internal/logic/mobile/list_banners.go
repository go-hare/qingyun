package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) ListBanners(ctx context.Context, request *mobile_fishing.ListBannersRequest, response *mobile_fishing.ListBannersResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "ListBanners",
	})

	ListBanners, err := models.ListBanners(func(session *xorm.Session) *xorm.Session {
		return session.Where("is_del = ?", 0).
			Where("banner_type = ?", request.BannerType).
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
		response.List = append(response.List, &mobile_fishing.Banner{
			ImageUrl:   ListBanners[i].ImageUrl,
			ViewType:   mobile_fishing.BannerViewType(ListBanners[i].ViewType),
			ExternalId: ListBanners[i].ExternalId,
			Title:      ListBanners[i].Title,
			LinkUrl:    ListBanners[i].LinkUrl,
			ThemeId:    ListBanners[i].ThemeId,
		})
	}
	return nil
}
