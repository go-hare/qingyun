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

func (m *MobileFishingService) UserArticleLike(ctx context.Context, request *mobile_fishing.UserArticleLikeRequest, response *mobile_fishing.UserArticleLikeResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "UserArticleLike",
	})
	if request.UserId < 1 || request.LikeId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	tableName := "user_like_" + strconv.Itoa(int(request.UserId)%models.UserLikeNum)
	GetUserLike, err := models.GetUserLike(func(session *xorm.Session) *xorm.Session {
		return session.Table(tableName).Where("user_like_type = ?", mobile_fishing.UserLikeType_user_like_type_active).
			Where("like_id = ?", request.LikeId).Where("user_id = ?", request.UserId)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUserLike",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	GetArticle, err := models.GetArticle(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.LikeId).
			Where("is_del = ?", 0).
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
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_user_like_error)
		return nil
	}

	if request.IfLike { //点赞
		if GetUserLike != nil {
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_user_like_error)
			return nil
		}
		if err = models.CreateUserLike(&models.UserLike{
			UserId:       request.UserId,
			LikeId:       request.LikeId,
			UserLikeType: int64(mobile_fishing.UserLikeType_user_like_type_active),
			LikeTime:     time.Now(),
		}); err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "CreateUserLike",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}

		models.UpdateArticle(func(session *xorm.Session) *xorm.Session {
			return session.Where("id = ?", request.LikeId).Incr("like_count")
		}, GetArticle)

	} else { //取消点赞
		if GetUserLike == nil {
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
			return nil
		}
		if err = models.DeleteUserLike(func(session *xorm.Session) *xorm.Session {
			return session.Table(tableName).Where("id = ?", GetUserLike.Id)
		}); err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "DeleteUserLike",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}

		models.UpdateArticle(func(session *xorm.Session) *xorm.Session {
			return session.Where("id = ?", request.LikeId).Decr("like_count")
		}, GetArticle)
	}
	return nil
}
