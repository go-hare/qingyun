package models

import "qingyun/common/store/mysql"

type ArticleCategory struct {
	mysql.Model `xorm:"extends"`
	Name        string `xorm:"not null default '' VARCHAR(256) comment('名称')" json:"name"`
	Weight      int64  `xorm:"not null default 0 INT(11) comment('权重')" json:"weight"`
	CreateTime  int64  `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
	IsDel       int64  `xorm:"not null default 0 INT(11) comment('是否删除 0 否 1是')" json:"is_del"`
	DelTime     int64  `xorm:"not null default 0 INT(11) comment('删除时间')" json:"del_time"`
}

func ListArticleCategorys(filter mysql.OrmFilter) (list []*ArticleCategory, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}
