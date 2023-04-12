package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) CheckInfomationOrderBuy(ctx context.Context, request *mobile_fishing.CheckInfomationOrderBuyRequest, response *mobile_fishing.CheckInfomationOrderBuyResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "CheckInfomationOrderBuy",
	})

	if request.InfomationId < 1 || request.UserId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetInfomationOrder, err := models.GetFishingOrder(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId).
			Where("goods_type = ?", mobile_fishing.FishingGoodsType_fishing_goods_type_infomotion).
			Where("goods_id = ?", request.InfomationId).Select("id")
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetFishingOrder",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	if GetInfomationOrder == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	response.IfBuy = true
	return nil

}
