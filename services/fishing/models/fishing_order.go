package models

import "qingyun/common/store/mysql"

type FishingOrder struct {
	mysql.Model `xorm:"extends"`
	Platform    int64 `xorm:"not null default 0 INT(11) comment('平台 0 小程序')" json:"platform"`
	CreateTime  int64 `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
	UserId      int64 `xorm:"not null default 0 INT(11) comment('用户id') index()" json:"user_id"`
	GoodsType   int64 `xorm:"not null default 0 INT(11) comment('类型 0 鱼获，1文章')" json:"goods_type"`
	GoodsId     int64 `xorm:"not null default 0 INT(11) comment('商品id')" json:"goods_id"`
	OrderPrice  int64 `xorm:"not null default 0 INT(11) comment('金额 分')" json:"order_price"`
	PayType     int64 `xorm:"not null default 0 INT(11) comment('支付类型 0 钓比 1积分')" json:"pay_type"`
	PayTime     int64 `xorm:"not null default 0 INT(11) comment('支付时间')" json:"pay_time"`
}

func GetFishingOrder(filter mysql.OrmFilter) (*FishingOrder, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(FishingOrder)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func CreateFishingOrder(message *FishingOrder) (err error) {
	_, err = mysql.GetDB().Insert(message)
	return
}

func ListFishingOrders(filter mysql.OrmFilter) (list []*FishingOrder, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}
