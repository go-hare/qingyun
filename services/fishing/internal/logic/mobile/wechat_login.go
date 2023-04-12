package mobile

import (
	"context"
	"github.com/realmicro/realmicro/common/util/time"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniconfig "github.com/silenceper/wechat/v2/miniprogram/config"
	log "github.com/sirupsen/logrus"
	"qingyun/common/token"
	"qingyun/services/fishing/models"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileFishingService) WechatLogin(ctx context.Context, request *mobile_fishing.WechatLoginRequest, response *mobile_fishing.WechatLoginResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "WechatLogin",
	})

	if request.Jscode == "" {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &miniconfig.Config{
		AppID:     "wx4eb85ed6f45d8c9b",
		AppSecret: "3c76c605b2749257c9ad01544c545500",
		Cache:     memory,
	}
	mini := wc.GetMiniProgram(cfg)
	a := mini.GetAuth()
	session, err := a.Code2Session(request.Jscode)
	if err != nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}
	encryptor := mini.GetEncryptor()
	wechatInfo, err := encryptor.Decrypt(session.SessionKey, request.EncryptedData, request.Iv)
	if err != nil {
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_param_error)
		return nil
	}

	openId := wechatInfo.OpenID
	//查询用户是否存在
	getUserPlatform, err := models.GetUserPlatform(func(session *xorm.Session) *xorm.Session {
		return session.Where("open_id = ?", openId).Where("platform_type = ?", mobile_fishing.ClientPlatformType_client_wechat_platform)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetUserPlatform",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
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
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}
		if mobileUser.UserStatus == int64(mobile_fishing.UserStatus_user_status_block) {
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_block_error)
			return nil
		}
	} else {
		CreateTime := time.Now()
		mobileUser = &models.User{
			Mobile:       wechatInfo.PhoneNumber,
			NickName:     wechatInfo.NickName,
			AvatarUrl:    wechatInfo.AvatarURL,
			City:         wechatInfo.City,
			Province:     wechatInfo.Province,
			Country:      wechatInfo.Country,
			Sex:          int64(wechatInfo.Gender),
			RegisterTime: CreateTime,
		}
		if err = models.CreateUser(mobileUser); err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "CreateUser",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}

		userPlatform := &models.UserPlatform{
			UserId:       mobileUser.Id,
			PlatformType: int64(mobile_fishing.ClientPlatformType_client_wechat_platform),
			OpenId:       openId,
			CreateTime:   CreateTime,
			UnionId:      session.UnionID,
		}
		if err = models.CreateUserPlatform(userPlatform); err != nil {
			logger.WithFields(log.Fields{
				"ErrorType": "Database",
				"Function":  "CreateUserPlatform",
			}).Error(err)
			mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
			return nil
		}
	}

	tokenStr, err := token.CreateJwtToken(mobileUser.Id)
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "token",
			"Function":  "CreateJwtToken",
		}).Error(err)
		mobile_fishing.ServiceStatus(response.Status, mobile_fishing.StatusCode_status_internal_error)
		return nil
	}
	response.Jwt = tokenStr
	response.Info = &mobile_fishing.User{
		UserId:       mobileUser.Id,
		Avatar:       mobileUser.AvatarUrl,
		Mobile:       mobileUser.Mobile,
		NickName:     mobileUser.NickName,
		RegisterTime: mobileUser.RegisterTime,
	}
	return nil
}
