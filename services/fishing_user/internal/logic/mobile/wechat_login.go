package mobile

import (
	"context"
	"github.com/realmicro/realmicro/common/util/time"
	log "github.com/sirupsen/logrus"
	"qingyun/common/token"
	"qingyun/services/fishing_user/models"
	mobile_fishing_user "qingyun/services/fishing_user/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingUserService) WechatLogin(ctx context.Context, request *mobile_fishing_user.WechatLoginRequest, response *mobile_fishing_user.WechatLoginResponse) (err error) {
	response.Status = &mobile_fishing_user.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "WechatLogin",
	})

	if request.Jscode == "" {
		mobile_fishing_user.ServiceStatus(response.Status, mobile_fishing_user.StatusCode_status_param_error)
		return nil
	}
	//GetSession, err := u.MiniProgram.GetSession(request.Jscode)
	//if err != nil {
	//	u.log.Error("dw")
	//	return nil, v12.ErrorUserNotFound("%s", "wwwww")
	//}
	//
	//return nil, status.Error(codes.Internal, "222")
	openId := "123434"
	//查询用户是否存在
	getUserPlatform, err := models.GetUserPlatform(func(session *xorm.Session) *xorm.Session {
		return session.Where("open_id = ?", openId).Where("platform_type = ?", mobile_fishing_user.ClientPlatformType_client_wechat_platform)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUserPlatform",
		}).Error(err)
		mobile_fishing_user.ServiceStatus(response.Status, mobile_fishing_user.StatusCode_status_internal_error)
	}
	var mobileUser *models.User
	if getUserPlatform != nil {
		mobileUser, err = models.GetUser(func(session *xorm.Session) *xorm.Session {
			return session.Where("id = ?", getUserPlatform.UserId)
		})
		if err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "GetUser",
			}).Error(err)
			mobile_fishing_user.ServiceStatus(response.Status, mobile_fishing_user.StatusCode_status_internal_error)
		}
	} else {
		CreateTime := time.Now()
		mobileUser = &models.User{
			Avatar:       "http://www.baidu.com",
			Mobile:       "",
			NickName:     "22222",
			AvatarUrl:    "",
			City:         "",
			Province:     "",
			Country:      "",
			Sex:          0,
			KeepClockDay: 0,
			BillDay:      0,
			BillNum:      0,
			RegisterTime: CreateTime,
		}
		if err = models.CreateUser(mobileUser); err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "CreateUser",
			}).Error(err)
			mobile_fishing_user.ServiceStatus(response.Status, mobile_fishing_user.StatusCode_status_internal_error)
		}

		userPlatform := &models.UserPlatform{
			UserId:       mobileUser.Id,
			PlatformType: int64(mobile_fishing_user.ClientPlatformType_client_wechat_platform),
			OpenId:       openId,
			CreateTime:   CreateTime,
		}
		if err = models.CreateUserPlatform(userPlatform); err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "CreateUserPlatform",
			}).Error(err)
			mobile_fishing_user.ServiceStatus(response.Status, mobile_fishing_user.StatusCode_status_internal_error)
		}
	}

	tokenStr, err := token.CreateJwtToken(mobileUser.Id)
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "token",
			"Function":  "CreateJwtToken",
		}).Error(err)
		mobile_fishing_user.ServiceStatus(response.Status, mobile_fishing_user.StatusCode_status_internal_error)
	}
	response.UserId = mobileUser.Id
	response.Jwt = tokenStr
	return nil
}
