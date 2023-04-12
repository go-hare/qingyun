package models

import "qingyun/common/store/mysql"

//收入明细
type UserCoinPayBill struct {
	mysql.Model `xorm:"extends"`
	UserId      int64 `xorm:"not null default 0 INT(11) comment('用户id')" json:"user_id"`
	CoinAmount  int64 `xorm:"not null default 0 INT(11) comment('钓币')" json:"coin_amount"`
	GoodsId     int64 `xorm:"not null default 0 INT(11) comment('商品id')" json:"goods_id"`
	GoodsType   int64 `xorm:"not null default 0 INT(11) comment('类型 0 鱼获 1文章')" json:"goods_type"`
	PayUserId   int64 `xorm:"not null default 0 INT(11) comment('购买用户')" json:"pay_user_id"`
	CreateTime  int64 `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
}
