package models

import "qingyun/common/store/mysql"

type Category struct { //分类
	mysql.Model    `xorm:"extends"`
	Pid            int64  `xorm:"not null default 0 INT(11) comment('父id')" json:"pid"`
	Name           string `xorm:"not null default '' VARCHAR(256) comment('名称')" json:"name"`
	Icon           string `xorm:"not null default '' VARCHAR(256) comment('头像')" json:"icon"`
	IsDelete       int64  `xorm:"not null default 0 INT(11) comment('是否删除 0否 1是')" json:"is_delete"`
	CategoryStatus int64  `xorm:"not null default 0 INT(11) comment('分类开启状态 0 否 1是')" json:"category_status"`
	Weight         int64  `xorm:"not null default 0 INT(11) comment('权重')" json:"weight"`
	CreateTime     int64  `xorm:"not null default 0 INT(11) comment('注册时间')" json:"create_time"`
}

func ListCategorys(filter mysql.OrmFilter) (list []*Category, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}
