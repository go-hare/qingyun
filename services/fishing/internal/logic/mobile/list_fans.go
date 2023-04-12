package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) ListFans(ctx context.Context, request *mobile_fishing.ListFansRequest, response *mobile_fishing.ListFansResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "ListFans",
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

	ListUserFollows, err := models.ListUserFollows(func(session *xorm.Session) *xorm.Session {
		return session.Where("follow_user_id = ?", request.UserId).
			OrderBy("follow_time desc").
			Limit(int(pageSize), int((page-1)*pageSize))
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListUserFollows",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	userIds := []int64{}
	for i := 0; i < len(ListUserFollows); i++ {
		userIds = append(userIds, ListUserFollows[i].UserId)
	}
	ListUsers, err := models.ListUsers(func(session *xorm.Session) *xorm.Session {
		return session.In("id", userIds)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListUsers",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	ListUserExtends, err := models.ListUserExtends(func(session *xorm.Session) *xorm.Session {
		return session.In("user_id", userIds)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListUserExtends",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	ListUserExtendsMap := make(map[int64]*models.UserExtend)
	for i := 0; i < len(ListUserExtends); i++ {
		ListUserExtendsMap[ListUserExtends[i].UserId] = ListUserExtends[i]
	}

	for i := 0; i < len(ListUsers); i++ {
		if v, ok := ListUserExtendsMap[ListUsers[i].Id]; ok {
			response.List = append(response.List, &mobile_fishing.User{
				UserId:   ListUsers[i].Id,
				Avatar:   ListUsers[i].AvatarUrl,
				NickName: ListUsers[i].NickName,
				UserExtend: &mobile_fishing.UserExtend{
					FollowCount:     0,
					FansCount:       0,
					InviteCount:     0,
					UnReadCount:     0,
					InfomationCount: v.InfomationCount,
					ArticleCount:    0,
				},
			})
		}

	}
	return nil
}
