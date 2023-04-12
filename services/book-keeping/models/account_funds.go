package models

import "qingyun/common/store/mysql"

type AccountFunds struct { //账户
	mysql.Model `xorm:"extends"`
	UserId      int64  `xorm:"not null default 0 INT(11) comment('用户id')" json:"user_id"`
	Name        string `xorm:"not null default '' VARCHAR(256) comment('名称')" json:"name"`
	Icon        string `xorm:"not null default '' VARCHAR(256) comment('图片')" json:"icon"`
	Weight      int64  `xorm:"not null default 0 INT(11) comment('权重')" json:"weight"`
	CreateTime  int64  `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
	IsDel       int64  `xorm:"not null default 0 INT(11) comment('是否删除')" json:"is_del"`
	DelTime     int64  `xorm:"not null default 0 INT(11) comment('删除时间')" json:"del_time"`
}

func ListAccountFunds(filter mysql.OrmFilter) (list []*AccountFunds, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func GetAccountFunds(filter mysql.OrmFilter) (accountFunds *AccountFunds, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(AccountFunds)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func CreateAccountFunds(bill *AccountFunds) (err error) {
	_, err = mysql.GetDB().Insert(bill)
	return
}

func UpdateAccountFunds(filter mysql.OrmFilter, message *AccountFunds) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Update(message)
	return
}
