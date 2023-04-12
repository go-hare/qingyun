package models

import "qingyun/common/store/mysql"

//按分类统计每年记账数据，不包括今年
type BillYearStat struct {
	mysql.Model `xorm:"extends"`
	UserId      int64 `xorm:"not null default 0 INT(11) comment('用户id')" json:"user_id"`
	CategoryId  int64 `xorm:"not null default 0 INT(11) comment('分类id')" json:"category_id"`
	Year        int64 `xorm:"not null default 0 INT(11) comment('年份')" json:"year"`
	Cost        int64 `xorm:"not null default 0 INT(11) comment('统计费用，单位分')" json:"cost"`
	LastCost    int64 `xorm:"not null default 0 INT(11) comment('上次统计，单位分')" json:"last_cost"`
	BillType    int64 `xorm:"not null default 0 INT(11) comment('类型 0 支出，1 收入')" json:"bill_type"`
}
