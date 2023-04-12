package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/book-keeping/models"
	mobile_book_keeping "qingyun/services/book-keeping/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileBookKeepingService) ListAccountFunds(ctx context.Context, request *mobile_book_keeping.ListAccountFundsRequest, response *mobile_book_keeping.ListAccountFundsResponse) (err error) {
	response.Status = &mobile_book_keeping.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "ListCategorys",
	})
	if request.UserId < 1 {
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_param_error)
		return nil
	}

	list, err := models.ListAccountFunds(func(session *xorm.Session) *xorm.Session {
		return session.Where("is_del = ?", 0).Where("user_id = ?", request.UserId)
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListAccountFunds",
		}).Error(err)
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_internal_error)
	}

	for i := 0; i < len(list); i++ {
		response.List = append(response.List, &mobile_book_keeping.AccountFunds{
			UserId:     list[i].UserId,
			Name:       list[i].Name,
			Icon:       list[i].Icon,
			Weight:     list[i].Weight,
			CreateTime: list[i].CreateTime,
		})
	}
	return nil
}
