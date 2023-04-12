package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/book-keeping/models"
	mobile_book_keeping "qingyun/services/book-keeping/proto/mobile"
	"time"
	"xorm.io/xorm"
)

func (m MobileBookKeepingService) CreateBill(ctx context.Context, request *mobile_book_keeping.CreateBillRequest, response *mobile_book_keeping.CreateBillResponse) (err error) {
	response.Status = &mobile_book_keeping.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "CreateBill",
	})
	if request.UserId < 1 || request.Bill == nil || request.Bill.AccountFundsId < 1 || request.Bill.CategoryId < 1 || request.Bill.CategoryId < 1 || request.Bill.Price < 1 {
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_param_error)
		return nil
	}
	GetAccountFunds, err := models.GetAccountFunds(func(session *xorm.Session) *xorm.Session {
		return session.Where("is_del = ?", 0).Where("user_id = ?", request.UserId).Where("id = ?", request.Bill.AccountFundsId)
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
	createTime := time.Now().Unix()
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	weekday := time.Now().Weekday()

	bill := &models.Bill{
		UserId:         request.UserId,
		Remark:         request.Bill.Remark,
		Price:          request.Bill.Price,
		BillType:       int64(request.Bill.BillType),
		AccountFundsId: request.Bill.AccountFundsId,
		Year:           int64(year),
		Week:           int64(weekday),
		Month:          int64(month),
		Day:            int64(day),
		CreateTime:     createTime,
	}
	if err = models.CreateBill(bill); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "CreateBill",
		}).Error(err)
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_internal_error)
	}
	return nil
}
