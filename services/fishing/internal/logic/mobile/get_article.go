package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"strconv"
	"xorm.io/xorm"
)

func (m *MobileFishingService) GetArticle(ctx context.Context, request *mobile_fishing.GetArticleRequest, response *mobile_fishing.GetArticleResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "GetArticle",
		"Method": "GetInfomation",
	})
	if request.ArticleId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetArticle, err := models.GetArticle(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.ArticleId).
			Where("article_status = ?", mobile_fishing.ActicleStatus_acticle_status_author_pass).Where("is_del = ?", 0)
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
	releaseUser, err := models.GetUser(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", GetArticle.UserId)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUser",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	if releaseUser == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	ListArticleTags, err := models.ListArticleTags(func(session *xorm.Session) *xorm.Session {
		return session.In("id", GetArticle.ListArticleTagIds)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListArticleTags",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	ListArticleTagsMap := make(map[int64]*mobile_fishing.ArticleTag)
	for i := 0; i < len(ListArticleTags); i++ {
		ListArticleTagsMap[ListArticleTags[i].Id] = &mobile_fishing.ArticleTag{
			TagId:   ListArticleTags[i].Id,
			TagName: ListArticleTags[i].Name,
		}
	}

	article := &mobile_fishing.Article{
		ArticleId:    GetArticle.Id,
		CategoryId:   GetArticle.CategoryId,
		UserId:       GetArticle.UserId,
		Title:        GetArticle.Title,
		CreateTime:   GetArticle.CreateTime,
		LikeCount:    GetArticle.LikeCount,
		CommentCount: GetArticle.CommentCount,
		CollectCount: GetArticle.CollectCount,
		ShareCount:   GetArticle.ShareCount,
		ViewCount:    GetArticle.ViewCount,
		ListTags:     GetArticle.ListTags,
		ImageUrl:     GetArticle.ImageUrl,
		Content:      GetArticle.Content,
		Avatar:       releaseUser.AvatarUrl,
		NickName:     releaseUser.NickName,
		IfLike:       false,
		IfCollection: false,
		ImageWidth:   GetArticle.ImageWidth,
		ImageHeight:  GetArticle.ImageHeight,
	}
	for j := 0; j < len(GetArticle.ListArticleTagIds); j++ {
		if v3, ok3 := ListArticleTagsMap[GetArticle.ListArticleTagIds[j]]; ok3 {
			article.ListArticleTags = append(article.ListArticleTags, v3)
		}
	}
	if request.UserId > 0 {
		//查询是否收藏是否点赞
		UserCollectiontableName := "user_collection_" + strconv.Itoa(int(request.UserId)%models.UserCollectionNum)
		GetUserCollection, err := models.GetUserCollection(func(session *xorm.Session) *xorm.Session {
			return session.Table(UserCollectiontableName).
				Where("user_id = ?", request.UserId).
				Where("user_collection_type = ?", mobile_fishing.UserColltionType_user_colltion_type_active).
				Where("collection_id = ?", request.ArticleId).Select("id")
		})
		if err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "GetUserCollection",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}
		if GetUserCollection != nil {
			article.IfCollection = true
		}
		tableName := "user_like_" + strconv.Itoa(int(request.UserId)%models.UserLikeNum)
		GetUserLike, err := models.GetUserLike(func(session *xorm.Session) *xorm.Session {
			return session.Table(tableName).
				Where("user_like_type = ?", mobile_fishing.UserLikeType_user_like_type_active).
				Where("like_id  = ?", request.ArticleId).Where("user_id = ?", request.UserId).Select("like_id")
		})
		if err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "GetUserLike",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}
		if GetUserLike != nil {
			article.IfLike = true
		}
	}
	response.Info = article
	return nil
}
