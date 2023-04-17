package models

import "qingyun/common/store/mysql"

type User struct {
	mysql.Model    `xorm:"extends"`
	UserName       string `json:"user_name"`
	OpenId         string `xorm:"not null default ''  VARCHAR(256) comment('第三方唯一标识') unique(user_id_open)" json:"open_id"`
	AvatarUrl      string `xorm:"not null default '' VARCHAR(256) comment('头像')" json:"avatar_url"`
	PointAmount    int64  `xorm:"not null default 0 INT(11) comment('积分')" json:"point_amount"`
	PayPointAmount int64  `xorm:"not null default 0 INT(11) comment('购买积分')" json:"pay_point_amount"`
	PointTotal     int64  `xorm:"not null default 0 INT(11) comment('累计积分')" json:"point_total"`
	CreateTime     int64  `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
}

func GetUser(filter mysql.OrmFilter) (user *User, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	user = &User{}
	_, err = session.Get(user)
	return
}

func CreateUser(user *User) (err error) {
	_, err = mysql.GetDB().Insert(user)
	return
}

func UpdateUser(filter mysql.OrmFilter, user *User) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Update(user)
	return
}
