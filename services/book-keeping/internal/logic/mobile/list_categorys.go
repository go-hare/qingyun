package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/book-keeping/models"
	mobile_book_keeping "qingyun/services/book-keeping/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileBookKeepingService) ListCategorys(ctx context.Context, request *mobile_book_keeping.ListCategorysRequest, response *mobile_book_keeping.ListCategorysResponse) (err error) {
	response.Status = &mobile_book_keeping.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "ListCategorys",
	})

	if request.Pid < 1 {
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_param_error)
		return nil
	}

	list, err := models.ListCategorys(func(session *xorm.Session) *xorm.Session {
		return session.Where("pid = ?", request.Pid).OrderBy("weight desc")
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListCategorys",
		}).Error(err)
		mobile_book_keeping.ServiceStatus(response.Status, mobile_book_keeping.StatusCode_status_internal_error)
	}

	for i := 0; i < len(list); i++ {
		response.List = append(response.List, &mobile_book_keeping.Category{
			Id:   list[i].Id,
			Name: list[i].Name,
			Icon: list[i].Icon,
		})
	}
	return nil
}
