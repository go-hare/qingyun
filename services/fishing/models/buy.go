package models

import "qingyun/common/store/mysql"

type Buy struct {
	mysql.Model   `xorm:"extends"`
	BuyStatus     int64 `xorm:"not null default 0 INT(11) comment('类型 0 鱼获，1文章')" json:"buy_status"`
	BuyExternalId int64 `xorm:"not null default 0 INT(11) comment('外部id') index(user_id)" json:"buy_external_id"`
	UserId        int64 `xorm:"not null default 0 INT(11) comment('用户id') index(user_id)" json:"user_id"`
	Price         int64 `xorm:"not null default 0 INT(11) comment('金额')" json:"price"`
	CreateTime    int64 `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
	PayTime       int64 `xorm:"not null default 0 INT(11) comment('支付时间')" json:"pay_time"`
}

func GetBuy(filter mysql.OrmFilter) (*Buy, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(Buy)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}
