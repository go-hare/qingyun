package models

import "qingyun/common/store/mysql"

type UserWechat struct {
	mysql.Model `xorm:"extends"`
	UserId      int64  `xorm:"not null default 0 INT(11) comment('用户id') index(user_id)" json:"user_id"`
	OpenId      string `xorm:"not null default '' CHAR(11) comment('用户唯一标识') unique" json:"open_id"`
	NickName    string `xorm:"not null default '' CHAR(11) comment('微信昵称')" json:"nick_name"`
	AvatarUrl   string `xorm:"not null default '' CHAR(11) comment('微信头像')" json:"avatar_url"`
	City        string `xorm:"not null default '' CHAR(11) comment('城市')" json:"city"`
	Province    string `xorm:"not null default '' CHAR(11) comment('省份')" json:"province"`
	Country     string `xorm:"not null default '' CHAR(11) comment('国家')" json:"country"`
	Gender      string `xorm:"not null default '' CHAR(11) comment('性别')" json:"gender"`
	CreateTime  int64  `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
}
