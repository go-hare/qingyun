package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"strconv"
	"xorm.io/xorm"
)

func (m *MobileFishingService) ListUserInfomations(ctx context.Context, request *mobile_fishing.ListUserInfomationsRequest, response *mobile_fishing.ListUserInfomationsResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "ListUserInfomations",
	})

	if request.PageSize < 1 || request.UserId < 1 || request.ReleaseUserId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	page := request.Page
	pageSize := request.PageSize
	if request.Page-1 <= 0 {
		page = 1
	}
	ListInfomations, err := models.ListInfomations(func(session *xorm.Session) *xorm.Session {
		return session.Where("infomation_status = ?", mobile_fishing.InfomationStatus_infomation_status_author_pass).
			Where("is_del = ?", 0).
			Where("user_id = ?", request.ReleaseUserId).
			OrderBy("create_time desc").Limit(int(pageSize), int((page-1)*pageSize))
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListInfomations",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	if len(ListInfomations) < 1 {
		return nil
	}
	userIds := []int64{}
	infomationIds := []int64{}
	for i := 0; i < len(ListInfomations); i++ {
		infomationIds = append(infomationIds, ListInfomations[i].Id)
		userIds = append(userIds, ListInfomations[i].UserId)
	}
	ListUsers, err := models.ListUsers(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_status = ?", mobile_fishing.UserStatus_user_status_default).In("id", userIds)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListUsers",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	ListUsersMap := make(map[int64]*models.User)
	for i := 0; i < len(ListUsers); i++ {
		ListUsersMap[ListUsers[i].Id] = ListUsers[i]
	}
	ListUserLikesMap := make(map[int64]struct{})
	ListUserCollectionsMap := make(map[int64]struct{})
	tableName := "user_like_" + strconv.Itoa(int(request.UserId)%models.UserLikeNum)
	ListUserLikes, err := models.ListUserLikes(func(session *xorm.Session) *xorm.Session {
		return session.Table(tableName).
			Where("user_like_type = ?", mobile_fishing.UserLikeType_user_like_type_infomation).
			In("like_id", infomationIds).Where("user_id = ?", request.UserId).Select("like_id")
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListUserLikes",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	for i := 0; i < len(ListUserLikes); i++ {
		ListUserLikesMap[ListUserLikes[i].LikeId] = struct{}{}
	}
	tableName = "user_collection_" + strconv.Itoa(int(request.UserId)%models.UserCollectionNum)
	ListUserCollections, err := models.ListUserCollections(func(session *xorm.Session) *xorm.Session {
		return session.Table(tableName).
			Where("user_collection_type = ?", mobile_fishing.UserColltionType_user_colltion_type_infomation).
			Where("user_collection_status = ?", 0).
			In("collection_id", infomationIds).Where("user_id = ?", request.UserId).Select("collection_id")
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListUserCollections",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	for i := 0; i < len(ListUserCollections); i++ {
		ListUserCollectionsMap[ListUserCollections[i].CollectionId] = struct{}{}
	}
	for i := 0; i < len(ListInfomations); i++ {
		if v, ok := ListUsersMap[ListInfomations[i].UserId]; ok {
			infomation := &mobile_fishing.Infomation{
				InfomationType:       mobile_fishing.InfomationType(ListInfomations[i].InfomationType),
				UserId:               ListInfomations[i].UserId,
				Country:              ListInfomations[i].Country,
				City:                 ListInfomations[i].City,
				District:             ListInfomations[i].District,
				Content:              ListInfomations[i].Content,
				Park:                 ListInfomations[i].Park,
				CreateTime:           ListInfomations[i].CreateTime,
				PayType:              mobile_fishing.PayType(ListInfomations[i].PayType),
				Price:                ListInfomations[i].Price,
				ListTags:             ListInfomations[i].ListTags,
				ListInfomationImages: ListInfomations[i].ListInfomationImages,
				LikeCount:            0,
				CommentCount:         0,
				CollectCount:         0,
				ReadCount:            0,
				Avatar:               v.AvatarUrl,
				NickName:             v.NickName,
				Distance:             10000,
				InfomationId:         ListInfomations[i].Id,
			}
			infomation.LikeCount = ListInfomations[i].LikeCount
			infomation.CommentCount = ListInfomations[i].CommentCount
			infomation.CollectCount = ListInfomations[i].CollectCount
			infomation.ReadCount = ListInfomations[i].ReadCount
			infomation.InfomationImageCount = int64(len(ListInfomations[i].ListInfomationImages))
			if _, ok2 := ListUserCollectionsMap[ListInfomations[i].Id]; ok2 {
				infomation.IfCollection = true
			}
			if _, ok2 := ListUserLikesMap[ListInfomations[i].Id]; ok2 {
				infomation.IfLike = true
			}
			response.List = append(response.List, infomation)
		}
	}
	return nil
}
