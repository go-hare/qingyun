package mobile

import (
	"context"
	"github.com/realmicro/realmicro/common/util/time"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) SendComment(ctx context.Context, request *mobile_fishing.SendCommentRequest, response *mobile_fishing.SendCommentResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "SendComment",
	})
	if request.UserId < 1 || request.ObjectId < 1 || request.Content == "" {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetUser, err := models.GetUser(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.UserId).
			Where("user_status = ?", mobile_fishing.UserStatus_user_status_default)
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUser",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if GetUser == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	if request.CommentType == mobile_fishing.CommentType_comment_type_infomation {

	}

	if request.CommentType == mobile_fishing.CommentType_comment_type_article {
		GetArticle, err := models.GetArticle(func(session *xorm.Session) *xorm.Session {
			return session.Where("id = ?", request.ObjectId).
				Where("article_status = ?", mobile_fishing.ActicleStatus_acticle_status_author_pass)
		})
		if err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "GetArticle",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}

		if GetArticle == nil {
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
			return nil
		}
	}

	comment := &models.Comment{
		CommentType:      int64(request.CommentType),
		CommentStatus:    0,
		CreateTime:       time.Now(),
		ObjectId:         request.ObjectId,
		ReleaseUserId:    request.UserId,
		LikeCount:        0,
		RepliedCount:     0,
		RepliedUserId:    0,
		RepliedCommentId: request.CommentId,
		Content:          request.Content,
	}
	if request.CommentId > 0 {
		GetComment, err := models.GetComment(func(session *xorm.Session) *xorm.Session {
			return session.Where("id = ?", request.CommentId)
		})
		if err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "GetArticle",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}
		if GetComment == nil {
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
			return nil
		}
		if GetComment.ReleaseUserId == request.UserId {
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
			return nil
		}
		comment.RepliedUserId = GetComment.ReleaseUserId
	}

	if err = models.CreateComment(comment); err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "CreateComment",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	if request.CommentType == mobile_fishing.CommentType_comment_type_article {
		models.UpdateArticle(func(session *xorm.Session) *xorm.Session {
			return session.Where("id = ?", request.ObjectId).Decr("comment_count")
		}, &models.Article{})
	} else {
		models.UpdateInfomation(func(session *xorm.Session) *xorm.Session {
			return session.Where("id = ?", request.ObjectId).Decr("comment_count")
		}, &models.Infomation{})
	}
	return nil

}
