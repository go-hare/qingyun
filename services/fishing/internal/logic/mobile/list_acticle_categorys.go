package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) ListActicleCategorys(ctx context.Context, request *mobile_fishing.ListActicleCategorysRequest, response *mobile_fishing.ListActicleCategorysResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "ListActicleCategorys",
	})

	ListArticleCategorys, err := models.ListArticleCategorys(func(session *xorm.Session) *xorm.Session {
		return session.Where("is_del = ?", 0).OrderBy("weight desc")
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListArticleCategorys",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	for i := 0; i < len(ListArticleCategorys); i++ {
		response.List = append(response.List, &mobile_fishing.ActicleCategory{
			CategoryId: ListArticleCategorys[i].Id,
			Name:       ListArticleCategorys[i].Name,
		})
	}
	return nil
}
