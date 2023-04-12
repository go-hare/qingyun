package mobile

import (
	"context"
	"errors"
	"github.com/realmicro/realmicro/common/util/time"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"qingyun/common/store/mysql"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) CreateInfomationOrder(ctx context.Context, request *mobile_fishing.CreateInfomationOrderRequest, response *mobile_fishing.CreateInfomationOrderResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "CreateInfomotionOrder",
	})

	if request.InfomationId < 1 || request.UserId < 1 {
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

	GetInfomation, err := models.GetInfomation(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.InfomationId).
			Where("is_del = ?", 0).
			Where("infomation_status = ?", mobile_fishing.InfomationStatus_infomation_status_author_pass)
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

	GetInfomationOrder, err := models.GetFishingOrder(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId).Where("goods_id = ?", request.InfomationId)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetFishingOrder",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	if GetInfomationOrder != nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	infomationOrder := &models.FishingOrder{
		UserId:     request.UserId,
		Platform:   0,
		GoodsId:    request.InfomationId,
		PayType:    int64(request.PayType),
		OrderPrice: GetInfomation.Price,
		GoodsType:  int64(mobile_fishing.FishingGoodsType_fishing_goods_type_infomotion),
		PayTime:    time.Now(),
		CreateTime: time.Now(),
	}
	if request.PayType == mobile_fishing.FishingOrderPayType_fishing_order_pay_type_coin { //钓币支付
		if GetUser.CoinAmount < GetInfomation.Price {
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_infomation_saldo_error)
			return nil
		}

		UserCoinPayBill := &models.UserCoinPayBill{
			UserId:     GetInfomation.UserId,
			CoinAmount: decimal.NewFromInt(GetInfomation.Price).Mul(decimal.NewFromFloat(0.8)).IntPart(),
			GoodsId:    GetInfomation.Id,
			GoodsType:  0,
			PayUserId:  GetUser.Id,
			CreateTime: time.Now(),
		}
		err = mysql.Transaction([]mysql.Handler{
			func(session *xorm.Session) error {
				affected, err := session.Where("id = ?", GetUser.Id).
					Where("coin_amount >= ?", GetInfomation.Price).
					Decr("coin_amount", GetInfomation.Price).Incr("coin_pay_total", GetInfomation.Price).Update(GetUser)
				if err != nil {
					return err
				}
				if affected == 0 {
					return errors.New("支付失败")
				}
				return nil
			},
			func(session *xorm.Session) error {
				_, err = session.Insert(UserCoinPayBill)
				return err
			},
			func(session *xorm.Session) error {
				_, err = session.Where("id = ?", GetInfomation.UserId).Incr("coin_total", UserCoinPayBill.CoinAmount).Update(&models.User{})
				return err
			},
		})

		//Affected, err := models.UpdateUserAffected(func(session *xorm.Session) *xorm.Session {
		//	return session.Where("id = ?", GetUser.Id).
		//		Where("coin_amount >= ?", GetInfomation.Price).
		//		Decr("coin_amount", GetInfomation.Price)
		//}, GetUser)
		if err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "Transaction",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_infomation_pay_error)
			return nil
		}
		//if Affected == 0 {
		//	mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_infomation_pay_error)
		//	return nil
		//}
	}
	if request.PayType == mobile_fishing.FishingOrderPayType_fishing_order_pay_type_point {
		if GetUser.PointAmount < GetInfomation.Price*10 {
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_infomation_saldo_error)
			return nil
		}
		Affected, err := models.UpdateUserAffected(func(session *xorm.Session) *xorm.Session {
			return session.Where("id = ?", GetUser.Id).
				Where("point_amount >= ?", GetInfomation.Price*10).
				Decr("point_amount", GetInfomation.Price*10)
		}, GetUser)
		if err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "UpdateUserAffected",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}
		if Affected == 0 {
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_infomation_pay_error)
			return nil
		}
	}
	if err = models.CreateFishingOrder(infomationOrder); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "CreateFishingOrder",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	response.Info = &mobile_fishing.InfomationOrder{
		OrderId:      infomationOrder.Id,
		UserId:       infomationOrder.UserId,
		InfomationId: infomationOrder.GoodsId,
		PayType:      mobile_fishing.FishingOrderPayType(infomationOrder.PayType),
		Price:        infomationOrder.OrderPrice,
		CreateTime:   infomationOrder.CreateTime,
	}
	return nil
}
