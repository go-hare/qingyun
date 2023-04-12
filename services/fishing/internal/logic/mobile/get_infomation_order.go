package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) GetInfomationOrder(ctx context.Context, request *mobile_fishing.GetInfomationOrderRequest, response *mobile_fishing.GetInfomationOrderResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "GetInfomationOrder",
	})

	if request.InfomationId < 1 || request.UserId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetFishingOrder, err := models.GetFishingOrder(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId).
			Where("goods_type = ?", mobile_fishing.FishingGoodsType_fishing_goods_type_infomotion).
			Where("goods_id = ?", request.InfomationId)
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetFishingOrder",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	if GetFishingOrder == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	response.OrderInfo = &mobile_fishing.InfomationOrder{
		OrderId:      GetFishingOrder.Id,
		UserId:       GetFishingOrder.UserId,
		InfomationId: GetFishingOrder.GoodsId,
		PayType:      mobile_fishing.FishingOrderPayType(GetFishingOrder.PayType),
		Price:        GetFishingOrder.OrderPrice,
		CreateTime:   GetFishingOrder.CreateTime,
	}

	GetInfomation, err := models.GetInfomation(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.InfomationId)
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
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	releaseUser, err := models.GetUser(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", GetInfomation.UserId)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUser",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	infomation := &mobile_fishing.Infomation{
		InfomationId:         GetInfomation.Id,
		InfomationType:       mobile_fishing.InfomationType(GetInfomation.InfomationType),
		UserId:               GetInfomation.UserId,
		Country:              GetInfomation.Country,
		City:                 GetInfomation.City,
		District:             GetInfomation.District,
		Content:              GetInfomation.Content,
		Park:                 GetInfomation.Park,
		CreateTime:           GetInfomation.CreateTime,
		PayType:              mobile_fishing.PayType(GetInfomation.PayType),
		Price:                GetInfomation.Price,
		ListTags:             GetInfomation.ListTags,
		ListInfomationImages: GetInfomation.ListInfomationImages,
		Avatar:               releaseUser.AvatarUrl,
		NickName:             releaseUser.NickName,
		LikeCount:            GetInfomation.LikeCount,
		CommentCount:         GetInfomation.CommentCount,
		CollectCount:         GetInfomation.CollectCount,
		ReadCount:            GetInfomation.ReadCount,
		Distance:             100000,
		InfomationImageCount: int64(len(GetInfomation.ListInfomationImages)),
		EnvImageCount:        int64(len(GetInfomation.ListEnvImages)),
		BaitImageCount:       int64(len(GetInfomation.ListBaits)),
		IfLike:               false,
		IfCollection:         false,
	}

	response.OrderInfo.Info = infomation
	return nil
}
