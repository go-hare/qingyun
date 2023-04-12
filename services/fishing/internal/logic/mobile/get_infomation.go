package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"strconv"
	"xorm.io/xorm"
)

func (m *MobileFishingService) GetInfomation(ctx context.Context, request *mobile_fishing.GetInfomationRequest, response *mobile_fishing.GetInfomationResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "GetInfomation",
	})

	if request.InfomationId < 1 {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	GetInfomation, err := models.GetInfomation(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.InfomationId).
			Where("is_del = ?", 0).
			Where("infomation_status = ?", mobile_fishing.InfomationStatus_infomation_status_author_pass)
	})

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetInfomation",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	if GetInfomation == nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	releaseUser, err := models.GetUser(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", GetInfomation.UserId)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUser",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	infomation := &mobile_fishing.Infomation{
		InfomationId:         GetInfomation.Id,
		InfomationType:       mobile_fishing.InfomationType(GetInfomation.InfomationType),
		UserId:               GetInfomation.UserId,
		Country:              GetInfomation.Country,
		City:                 GetInfomation.City,
		District:             GetInfomation.District,
		Content:              GetInfomation.Content,
		Park:                 GetInfomation.Park,
		CreateTime:           GetInfomation.CreateTime,
		PayType:              mobile_fishing.PayType(GetInfomation.PayType),
		Price:                GetInfomation.Price,
		ListTags:             GetInfomation.ListTags,
		ListInfomationImages: GetInfomation.ListInfomationImages,
		Avatar:               releaseUser.AvatarUrl,
		NickName:             releaseUser.NickName,
		LikeCount:            GetInfomation.LikeCount,
		CommentCount:         GetInfomation.CommentCount,
		CollectCount:         GetInfomation.CollectCount,
		ReadCount:            GetInfomation.ReadCount,
		Distance:             100000,
		InfomationImageCount: int64(len(GetInfomation.ListInfomationImages)),
		EnvImageCount:        int64(len(GetInfomation.ListEnvImages)),
		BaitImageCount:       int64(len(GetInfomation.ListBaits)),
		IfLike:               false,
		IfCollection:         false,
		Latitude:             request.Latitude,
		Longitude:            request.Longitude,
	}

	if request.UserId > 0 {
		//查询是否收藏是否点赞
		UserCollectiontableName := "user_collection_" + strconv.Itoa(int(request.UserId)%models.UserCollectionNum)
		GetUserCollection, err := models.GetUserCollection(func(session *xorm.Session) *xorm.Session {
			return session.Table(UserCollectiontableName).
				Where("user_id = ?", request.UserId).
				Where("user_collection_type = ?", mobile_fishing.UserColltionType_user_colltion_type_infomation).
				Where("collection_id = ?", request.InfomationId).Select("id")
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
			infomation.IfCollection = true
		}
		tableName := "user_like_" + strconv.Itoa(int(request.UserId)%models.UserLikeNum)
		GetUserLike, err := models.GetUserLike(func(session *xorm.Session) *xorm.Session {
			return session.Table(tableName).
				Where("user_like_type = ?", mobile_fishing.UserLikeType_user_like_type_infomation).
				Where("like_id  = ?", request.InfomationId).Where("user_id = ?", request.UserId).Select("like_id")
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
			infomation.IfLike = true
		}
		GetBuy, err := models.GetFishingOrder(func(session *xorm.Session) *xorm.Session {
			return session.Where("user_id = ?", request.UserId).
				Where("goods_type = ?", mobile_fishing.FishingGoodsType_fishing_goods_type_infomotion).
				Where("goods_id = ?", request.InfomationId).Select("id")
		})
		if err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "GetFishingOrder",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}
		if GetBuy != nil {
			infomation.IfBuy = true
		}
		GetUserFollow, err := models.GetUserFollow(func(session *xorm.Session) *xorm.Session {
			return session.Where("user_id = ?", request.UserId).Where("follow_user_id = ?", GetInfomation.UserId).Select("id")
		})
		if err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "GetUserFollow",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}
		if GetUserFollow != nil {
			infomation.IsFollow = true
		}

	}
	if infomation.IfBuy {
		infomation.ListEnvImages = GetInfomation.ListEnvImages
		ListBaits := []*mobile_fishing.Bait{}
		for i := 0; i < len(GetInfomation.ListBaits); i++ {
			ListBaits = append(ListBaits, &mobile_fishing.Bait{
				Image:       GetInfomation.ListBaits[i].Image,
				BaitName:    GetInfomation.ListBaits[i].BaitName,
				BaitPercent: GetInfomation.ListBaits[i].BaitPercent,
			})
		}
		infomation.ListBaits = ListBaits
	}
	models.UpdateInfomation(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", GetInfomation.Id).Incr("read_count")
	}, GetInfomation)
	response.Info = infomation
	return nil
}
