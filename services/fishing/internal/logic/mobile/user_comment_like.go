package mobile

import (
	"context"
	"github.com/realmicro/realmicro/common/util/time"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"strconv"
	"xorm.io/xorm"
)

func (m *MobileFishingService) UserCommentLike(ctx context.Context, request *mobile_fishing.UserCommentLikeRequest, response *mobile_fishing.UserCommentLikeResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "UserArticleLike",
	})
	if request.UserId < 1 || request.CommentId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetComment, err := models.GetComment(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.CommentId).
			Where("comment_type = ?", request.CommentType).
			Where("comment_status = ?", mobile_fishing.CommentType_comment_type_article)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetComment",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if GetComment == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	tableName := "user_comment_like_" + strconv.Itoa(int(request.UserId)%models.UserCommentLikeNum)
	GetUserCommentLike, err := models.GetUserCommentLike(func(session *xorm.Session) *xorm.Session {
		return session.Table(tableName).
			Where("comment_id = ?", request.CommentId).
			Where("user_id = ?", request.UserId)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUserLike",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	if request.IfLike {
		if GetUserCommentLike != nil {
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
			return nil
		}
		UserCommentLike := &models.UserCommentLike{
			UserId:    request.UserId,
			CommentId: request.CommentId,
			LikeTime:  time.Now(),
		}

		if err = models.CreateUserCommentLike(UserCommentLike); err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "CreateUserCommentLike",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}
		models.UpdateComment(func(session *xorm.Session) *xorm.Session {
			return session.Where("id = ?", request.CommentId).Incr("like_count")
		}, GetComment)

	} else {
		if GetUserCommentLike == nil {
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
			return nil
		}
		if err = models.DeleteUserCommentLike(func(session *xorm.Session) *xorm.Session {
			return session.Table(tableName).Where("id = ?", GetUserCommentLike.Id)
		}); err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "DeleteUserCommentLike",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}

		models.UpdateComment(func(session *xorm.Session) *xorm.Session {
			return session.Where("id = ?", request.CommentId).Decr("like_count")
		}, GetComment)
	}
	return nil
}
