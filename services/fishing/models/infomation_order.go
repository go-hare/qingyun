package models

import "qingyun/common/store/mysql"

type InfomationOrder struct {
	mysql.Model  `xorm:"extends"`
	UserId       int64 `xorm:"not null default 0 INT(11) comment('用户id') unique(user_infomation)" json:"user_id"`
	Platform     int64 `xorm:"not null default 0 INT(11) comment('平台 0 小程序')" json:"platform"`
	InfomationId int64 `xorm:"not null default 0 INT(11) comment('信息id') unique(user_infomation)" json:"infomation_id"`
	PayType      int64 `xorm:"not null default 0 INT(11) comment('支付类型')" json:"pay_type"`
	Price        int64 `xorm:"not null default 0 INT(11) comment('价格')" json:"price"`
	PayTime      int64 `xorm:"not null default 0 INT(11) comment('支付时间')" json:"pay_time"`
	CreateTime   int64 `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
}

func GetInfomationOrder(filter mysql.OrmFilter) (*InfomationOrder, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(InfomationOrder)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func CreateInfomationOrder(message *InfomationOrder) (err error) {
	_, err = mysql.GetDB().Insert(message)
	return
}

func ListInfomationOrders(filter mysql.OrmFilter) (list []*InfomationOrder, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}
