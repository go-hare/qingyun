package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"strconv"
	"xorm.io/xorm"
)

func (m *MobileFishingService) ListMessages(ctx context.Context, request *mobile_fishing.ListMessagesRequest, response *mobile_fishing.ListMessagesResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "ListMessages",
	})

	if request.PageSize < 1 || request.UserId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	page := request.Page
	pageSize := request.PageSize
	if request.Page-1 <= 0 {
		page = 1
	}
	tableName := "message_" + strconv.Itoa(int(request.UserId)%models.MessageNum)

	ListMessages, err := models.ListMessages(func(session *xorm.Session) *xorm.Session {
		return session.Table(tableName).
			Where("user_id = ?", request.UserId).OrderBy("read_status asc").
			OrderBy("create_time desc").Limit(int(pageSize), int((page-1)*pageSize))
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListMessages",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	for i := 0; i < len(ListMessages); i++ {
		response.List = append(response.List, &mobile_fishing.Message{
			MessageId:  ListMessages[i].Id,
			CreateTime: ListMessages[i].CreateTime,
			ReadStatus: mobile_fishing.MessageReadStatus(ListMessages[i].ReadStatus),
			Title:      ListMessages[i].Title,
			Content:    ListMessages[i].Content,
			Avatar:     ListMessages[i].Avatar,
		})
	}
	return nil
}
