package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) ListUserInfomationOrders(ctx context.Context, request *mobile_fishing.ListUserInfomationOrdersRequest, response *mobile_fishing.ListUserInfomationOrdersResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "ListUserInfomationOrders",
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

	ListInfomationOrders, err := models.ListFishingOrders(func(session *xorm.Session) *xorm.Session {
		return session.Where("user_id = ?", request.UserId).
			Where("goods_type = ?", mobile_fishing.FishingGoodsType_fishing_goods_type_infomotion).
			OrderBy("create_time desc").Limit(int(pageSize), int((page-1)*pageSize))
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListInfomationOrders",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}

	if len(ListInfomationOrders) < 1 {
		return nil
	}
	InfomationOrderIds := []int64{}
	ListInfomationOrdersMap := make(map[int64]*models.FishingOrder)
	for i := 0; i < len(ListInfomationOrders); i++ {
		InfomationOrderIds = append(InfomationOrderIds, ListInfomationOrders[i].GoodsId)
		ListInfomationOrdersMap[ListInfomationOrders[i].GoodsId] = ListInfomationOrders[i]
	}

	ListInfomations, err := models.ListInfomations(func(session *xorm.Session) *xorm.Session {
		return session.In("id", InfomationOrderIds)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListInfomations",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
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
				InfomationId:         ListInfomations[i].Id,
			}

			if v1, ok1 := ListInfomationOrdersMap[ListInfomations[i].Id]; ok1 {
				infomationOrder := &mobile_fishing.InfomationOrder{
					OrderId:      v1.Id,
					UserId:       v1.UserId,
					InfomationId: v1.GoodsId,
					PayType:      mobile_fishing.FishingOrderPayType(v1.PayType),
					Price:        v1.OrderPrice,
					CreateTime:   v1.CreateTime,
					Info:         infomation,
				}
				response.List = append(response.List, infomationOrder)
			}
		}

	}

	return nil
}
