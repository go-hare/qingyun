package models

import "qingyun/common/store/mysql"

type WechatBot struct {
	mysql.Model  `xorm:"extends"`
	AppKey       string `xorm:"not null default '' VARCHAR(128) comment('用户id')" json:"app_key"`
	UserName     string `xorm:"not null default '' VARCHAR(128) comment('名字')" json:"user_name"`
	NickName     string `xorm:"not null default '' VARCHAR(50) comment('昵称')" json:"nick_name"`
	HeadImgUrl   string `xorm:"not null default '' VARCHAR(128) comment('微信头像')" json:"avatar_url"`
	Session      string `xorm:"not null  TEXT json comment('图片列表')"  json:"session"`
	IfLogin      int64  `xorm:"not null default 0 INT(11) comment('是否登录')" json:"if_login"`
	CreateTime   int64  `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
	LoginTime    int64  `xorm:"not null default 0 INT(11) comment('登录时间')" json:"login_time"`
	LoginOutTime int64  `xorm:"not null default 0 INT(11) comment('退出时间')" json:"login_out_time"`
}

func GetWechatBot(filter mysql.OrmFilter) (*WechatBot, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(WechatBot)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}
func CreateWechatBot(message *WechatBot) (err error) {
	_, err = mysql.GetDB().Insert(message)
	return
}
func UpdateWechatBot(filter mysql.OrmFilter, message *WechatBot) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Update(message)
	return
}

func ListWechatBots(filter mysql.OrmFilter) (list []*WechatBot, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}
