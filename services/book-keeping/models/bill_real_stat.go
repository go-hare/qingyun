package models

import "qingyun/common/store/mysql"

//按分类统计本周本月今年的记账数据
type BillRealStat struct {
	mysql.Model `xorm:"extends"`
	UserId      int64 `xorm:"not null default 0 INT(11) comment('用户id')" json:"user_id"`
	CategoryId  int64 `xorm:"not null default 0 INT(11) comment('分类id')" json:"category_id"`
	Week        int64 `xorm:"not null default 0 INT(11) comment('当前周')" json:"week"`
	WeekCost    int64 `xorm:"not null default 0 INT(11) comment('当前周费用，单位分')" json:"week_cost"`
	Month       int64 `xorm:"not null default 0 INT(11) comment('当前月')" json:"month"`
	MonthCost   int64 `xorm:"not null default 0 INT(11) comment('当前月费用，单位分')" json:"month_cost"`
	Year        int64 `xorm:"not null default 0 INT(11) comment('当前年')" json:"year"`
	YearCost    int64 `xorm:"not null default 0 INT(11) comment('当前年费用，单位分')" json:"year_cost"`
	BillType    int64 `xorm:"not null default 0 INT(11) comment('类型 0 支出，1 收入')" json:"bill_type"`
}
