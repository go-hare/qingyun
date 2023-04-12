package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"strconv"
	"xorm.io/xorm"
)

func (m *MobileFishingService) ListActiclesByCategory(ctx context.Context, request *mobile_fishing.ListActiclesByCategoryRequest, response *mobile_fishing.ListActiclesByCategoryResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "ListActiclesByCategory",
	})

	if request.CategoryId < 1 || request.PageSize < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	page := request.Page
	pageSize := request.PageSize
	if request.Page-1 <= 0 {
		page = 1
	}
	ListArticles, err := models.ListArticles(func(session *xorm.Session) *xorm.Session {
		if request.CategoryId > 1 {
			session = session.Where("category_id = ?", request.CategoryId)
		}
		return session.Where("article_status = ?", mobile_fishing.ActicleStatus_acticle_status_author_pass).
			OrderBy("view_count desc").Where("is_del = ?", 0).
			OrderBy("create_time desc").Limit(int(pageSize), int((page-1)*pageSize))
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListArticles",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if len(ListArticles) < 1 {
		return nil
	}
	userIds := []int64{}
	acticleIds := []int64{}
	for i := 0; i < len(ListArticles); i++ {
		acticleIds = append(acticleIds, ListArticles[i].Id)
		userIds = append(userIds, ListArticles[i].UserId)
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
	if request.UserId > 0 {
		tableName := "user_like_" + strconv.Itoa(int(request.UserId)%models.UserLikeNum)
		ListUserLikes, err := models.ListUserLikes(func(session *xorm.Session) *xorm.Session {
			return session.Table(tableName).
				Where("user_like_type = ?", mobile_fishing.UserLikeType_user_like_type_active).
				In("like_id", acticleIds).Where("user_id = ?", request.UserId).Select("like_id")
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
				Where("user_collection_type = ?", mobile_fishing.UserColltionType_user_colltion_type_active).
				Where("user_collection_status = ?", 0).
				In("collection_id", acticleIds).Where("user_id = ?", request.UserId).Select("collection_id")
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
	}

	tagIds := []int64{}
	for i := 0; i < len(ListArticles); i++ {
		tagIds = append(tagIds, ListArticles[i].ListArticleTagIds...)
	}

	ListArticleTags, err := models.ListArticleTags(func(session *xorm.Session) *xorm.Session {
		return session.In("id", tagIds)
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
	for i := 0; i < len(ListArticles); i++ {
		if v, ok := ListUsersMap[ListArticles[i].UserId]; ok {
			acticle := &mobile_fishing.Article{
				ArticleId:    ListArticles[i].Id,
				CategoryId:   ListArticles[i].CategoryId,
				UserId:       ListArticles[i].UserId,
				Title:        ListArticles[i].Title,
				CreateTime:   ListArticles[i].CreateTime,
				LikeCount:    ListArticles[i].LikeCount,
				CommentCount: ListArticles[i].CommentCount,
				CollectCount: ListArticles[i].CollectCount,
				ShareCount:   ListArticles[i].ShareCount,
				ViewCount:    ListArticles[i].ViewCount,
				ListTags:     ListArticles[i].ListTags,
				ImageUrl:     ListArticles[i].ImageUrl,
				ImageHeight:  ListArticles[i].ImageHeight,
				ImageWidth:   ListArticles[i].ImageWidth,
				Avatar:       v.AvatarUrl,
				NickName:     v.NickName,
			}
			for j := 0; j < len(ListArticles[i].ListArticleTagIds); j++ {
				if v3, ok3 := ListArticleTagsMap[ListArticles[i].ListArticleTagIds[j]]; ok3 {
					acticle.ListArticleTags = append(acticle.ListArticleTags, v3)
				}
			}
			if _, ok2 := ListUserCollectionsMap[ListArticles[i].Id]; ok2 {
				acticle.IfCollection = true
			}
			if _, ok2 := ListUserLikesMap[ListArticles[i].Id]; ok2 {
				acticle.IfLike = true
			}
			response.List = append(response.List, acticle)
		}
	}
	return nil

}
