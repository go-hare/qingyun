package mobile

import (
	"context"
	"github.com/realmicro/realmicro/common/util/time"
	log "github.com/sirupsen/logrus"
	"qingyun/services/book-keeping/models"
	mobile_book_keeping "qingyun/services/book-keeping/proto/mobile"
)

func (m *MobileBookKeepingService) CreatetAccountFunds(ctx context.Context, request *mobile_book_keeping.CreatetAccountFundsRequest, response *mobile_book_keeping.CreatetAccountFundsResponse) (err error) {
	response.Status = &mobile_book_keeping.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "CreatetAccountFunds",
	})
	if request.UserId < 1 || request.Info == nil || request.Info.Name == "" {
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_param_error)
		return nil
	}

	accountFunds := &models.AccountFunds{
		UserId:     request.UserId,
		Name:       request.Info.Name,
		Icon:       request.Info.Icon,
		Weight:     request.Info.Weight,
		CreateTime: time.Now(),
	}

	if err = models.CreateAccountFunds(accountFunds); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "CreateAccountFunds",
		}).Error(err)
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_internal_error)
	}
	return nil
}
