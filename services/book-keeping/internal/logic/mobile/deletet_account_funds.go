package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/book-keeping/models"
	mobile_book_keeping "qingyun/services/book-keeping/proto/mobile"
	"time"
	"xorm.io/xorm"
)

func (m *MobileBookKeepingService) DeletetAccountFunds(ctx context.Context, request *mobile_book_keeping.DeletetAccountFundsRequest, response *mobile_book_keeping.DeletetAccountFundsResponse) (err error) {
	response.Status = &mobile_book_keeping.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "DeletetAccountFunds",
	})
	if request.UserId < 1 || request.AccountFundsId < 1 {
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_param_error)
		return nil
	}

	GetAccountFunds, err := models.GetAccountFunds(func(session *xorm.Session) *xorm.Session {
		return session.Where("is_del = ?", 0).Where("user_id = ?", request.UserId).Where("id = ?", request.AccountFundsId)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetAccountFunds",
		}).Error(err)
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_internal_error)
	}
	if GetAccountFunds == nil {
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_param_error)
		return nil
	}
	GetBill, err := models.GetBill(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId).Where("account_funds_id = ?", request.AccountFundsId)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetAccountFunds",
		}).Error(err)
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_internal_error)
	}
	if GetBill != nil {
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_has_bill_error)
		return nil
	}
	err = models.UpdateAccountFunds(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", GetAccountFunds.Id).Cols("is_del,del_time")
	}, &models.AccountFunds{IsDel: 1, DelTime: time.Now().Unix()})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "UpdateAccountFunds",
		}).Error(err)
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_internal_error)
	}
	return nil
}
