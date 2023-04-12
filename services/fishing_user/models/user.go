package models

import "qingyun/common/store/mysql"

type User struct {
	mysql.Model  `xorm:"extends"`
	Avatar       string `xorm:"not null default '' VARCHAR(256) comment('头像')" json:"avatar"`
	Mobile       string `xorm:"not null default '' CHAR(11) comment('电话')" json:"mobile"`
	NickName     string `xorm:"not null default '' CHAR(11) comment('微信昵称')" json:"nick_name"`
	AvatarUrl    string `xorm:"not null default '' CHAR(11) comment('微信头像')" json:"avatar_url"`
	City         string `xorm:"not null default '' CHAR(11) comment('城市')" json:"city"`
	Province     string `xorm:"not null default '' CHAR(11) comment('省份')" json:"province"`
	Country      string `xorm:"not null default '' CHAR(11) comment('国家')" json:"country"`
	Sex          int64  `xorm:"not null default 0 INT(11) comment('性别：1=男，2=女，3=保密')" json:"sex"`
	KeepClockDay int64  `xorm:"not null default 0 INT(11) comment('连续打卡天数')" json:"keep_clock_day"`
	BillDay      int64  `xorm:"not null default 0 INT(11) comment('记账天数')" json:"bill_day"`
	BillNum      int64  `xorm:"not null default 0 INT(11) comment('总记账次数')" json:"bill_num"`
	RegisterTime int64  `xorm:"not null default 0 INT(11) comment('注册时间')" json:"register_time"`
}

func GetUser(filter mysql.OrmFilter) (*User, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(User)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func ListUsers(filter mysql.OrmFilter) (list []*User, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func ListUsersCount(filter mysql.OrmFilter) (list []*User, count int64, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}

	list = make([]*User, 0)
	count, err = session.FindAndCount(&list)
	return
}

func CreateUser(message *User) (err error) {
	_, err = mysql.GetDB().Insert(message)
	return
}

func UpdateUser(filter mysql.OrmFilter, message *User) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Update(message)
	return
}

func DeleteUser(filter mysql.OrmFilter) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Delete(&User{})
	return
}
