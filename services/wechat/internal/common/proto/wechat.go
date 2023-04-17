package proto

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"qingyun/services/wechat/models"
	"xorm.io/xorm"
)

type MysqlHotReloadStorage struct {
	appKey string
	reader *bytes.Reader
}

func NewMysqlHotReloadStorage(appKey string) *MysqlHotReloadStorage {
	return &MysqlHotReloadStorage{appKey: appKey}
}

// Load 重写热登录数据加载，从Redis取数据
func (f *MysqlHotReloadStorage) Read(p []byte) (n int, err error) {
	if f.reader == nil {
		// 从Redis获取热登录数据
		data, err := models.GetWechatBot(func(session *xorm.Session) *xorm.Session {
			return session.Where("app_key = ?", f.appKey)
		})
		if err != nil {
			log.Errorf("读取热登录数据失败: %v", err)
			return 0, err
		}
		f.reader = bytes.NewReader([]byte(data.Session))
	}
	return f.reader.Read(p)
}

// Dump 重写更新热登录数据，保存到Redis
func (f *MysqlHotReloadStorage) Write(p []byte) (n int, err error) {
	GetWechatBot, err := models.GetWechatBot(func(session *xorm.Session) *xorm.Session {
		return session.Where("app_key = ?", f.appKey)
	})
	if err != nil {
		log.Errorf("保存微信热登录信息失败: %v", err.Error())
		return
	}
	if GetWechatBot == nil {
		models.CreateWechatBot(&models.WechatBot{
			AppKey: f.appKey,
		})

	}
	err = models.UpdateWechatBot(func(session *xorm.Session) *xorm.Session {
		return session.Where("app_key = ?", f.appKey).Cols("session")
	}, &models.WechatBot{Session: string(p)})
	if err != nil {
		log.Errorf("保存微信热登录信息失败: %v", err.Error())
		return
	}
	return len(p), nil
}

// Close 需要关闭
func (f *MysqlHotReloadStorage) Close() error {
	f.reader = nil
	return nil
}
