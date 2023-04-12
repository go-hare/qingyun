package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/book-keeping/models"
	mobile_book_keeping "qingyun/services/book-keeping/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileBookKeepingService) TransferBillToAccountFunds(ctx context.Context, request *mobile_book_keeping.TransferBillToAccountFundsRequest, response *mobile_book_keeping.TransferBillToAccountFundsResponse) (err error) {
	response.Status = &mobile_book_keeping.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "TransferBillToAccountFunds",
	})
	if request.UserId < 1 || request.BillId < 1 || request.AccountFundsId < 1 {
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
		return session.Where("user_id = ?", request.UserId).Where("id = ?", request.BillId)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetBill",
		}).Error(err)
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_internal_error)
	}
	if GetBill == nil {
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_param_error)
		return nil
	}
	GetBill.AccountFundsId = GetAccountFunds.Id
	if err = models.UpdateBill(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", GetBill.Id).Cols("account_funds_id")
	}, GetBill); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "UpdateBill",
		}).Error(err)
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_internal_error)
	}
	return nil
}
