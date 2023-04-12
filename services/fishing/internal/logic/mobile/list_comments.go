package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"strconv"
	"xorm.io/xorm"
)

func (m *MobileFishingService) ListComments(ctx context.Context, request *mobile_fishing.ListCommentsRequest, response *mobile_fishing.ListCommentsResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "ListComments",
	})

	if request.PageSize < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	page := request.Page
	pageSize := request.PageSize
	if request.Page-1 <= 0 {
		page = 1
	}

	ListComments, err := models.ListComments(func(session *xorm.Session) *xorm.Session {
		return session.Where("comment_type = ?", request.CommentType).
			Where("comment_status = ?", mobile_fishing.CommentStatus_comment_status_author_pass).
			Where("object_id = ?", request.ObjectId).
			OrderBy("create_time desc").Limit(int(pageSize), int((page-1)*pageSize))
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListComments",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if len(ListComments) < 1 {
		return nil
	}

	userIds := []int64{}
	commentIds := []int64{}
	for i := 0; i < len(ListComments); i++ {
		commentIds = append(commentIds, ListComments[i].Id)
		userIds = append(userIds, ListComments[i].ReleaseUserId)
		if ListComments[i].RepliedUserId > 0 {
			userIds = append(userIds, ListComments[i].RepliedUserId)
		}
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
	ListUsersMap := make(map[int64]*models.User)
	for i := 0; i < len(ListUsers); i++ {
		ListUsersMap[ListUsers[i].Id] = ListUsers[i]
	}
	ListUserCommentLikesMap := make(map[int64]struct{})
	if request.UserId > 0 {
		tableName := "user_comment_like_" + strconv.Itoa(int(request.UserId)%models.UserCommentLikeNum)
		ListUserCommentLikes, err := models.ListUserCommentLikes(func(session *xorm.Session) *xorm.Session {
			return session.Table(tableName).
				In("comment_id", commentIds).Where("user_id = ?", request.UserId).Select("comment_id")
		})
		if err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "ListUserCommentLikes",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}
		for i := 0; i < len(ListUserCommentLikes); i++ {
			ListUserCommentLikesMap[ListUserCommentLikes[i].CommentId] = struct{}{}
		}
	}

	for i := 0; i < len(ListComments); i++ {
		comment := &mobile_fishing.Comment{
			CommentId:        ListComments[i].Id,
			CreateTime:       ListComments[i].CreateTime,
			ObjectId:         ListComments[i].ObjectId,
			ReleaseUserId:    ListComments[i].ReleaseUserId,
			LikeCount:        ListComments[i].LikeCount,
			RepliedCount:     ListComments[i].RepliedCount,
			RepliedUserId:    ListComments[i].RepliedUserId,
			RepliedCommentId: ListComments[i].RepliedCommentId,
			Content:          ListComments[i].Content,
		}
		if _, ok := ListUserCommentLikesMap[ListComments[i].Id]; ok {
			comment.IfLike = true
		}
		if v1, ok1 := ListUsersMap[ListComments[i].ReleaseUserId]; ok1 {
			comment.Avatar = v1.AvatarUrl
			comment.NickName = v1.NickName
		}
		if v2, ok2 := ListUsersMap[ListComments[i].RepliedUserId]; ok2 {
			comment.Content = v2.NickName + ": " + comment.Content
		}
		response.List = append(response.List, comment)
	}
	return nil
}
