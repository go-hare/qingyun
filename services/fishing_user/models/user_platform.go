package models

import "qingyun/common/store/mysql"

type UserPlatform struct {
	mysql.Model  `xorm:"extends"`
	UserId       int64  `xorm:"not null default 0 INT(11) comment('用户id') index(user_id)" json:"user_id"`
	PlatformType int64  `xorm:"not null default 0 INT(11) comment('第三方平台类型 0 微信')" json:"platform_type"`
	OpenId       string `xorm:"not null default '' CHAR(11) comment('第三方唯一标识')" json:"open_id"`
	CreateTime   int64  `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
}

func GetUserPlatform(filter mysql.OrmFilter) (*UserPlatform, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(UserPlatform)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func CreateUserPlatform(message *UserPlatform) (err error) {
	_, err = mysql.GetDB().Insert(message)
	return
}
