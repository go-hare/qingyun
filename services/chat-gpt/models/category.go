package models

import "qingyun/common/store/mysql"

//分类
type Category struct {
	mysql.Model `xorm:"extends"`
	Name        string `xorm:"not null default '' VARCHAR(256) comment('标题')" json:"name"`
	Desc        string `xorm:"not null default '' VARCHAR(256) comment('描述')" json:"desc"`
	CreateTime  int64  `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
	Wgiht       int64  `xorm:"not null default 0 INT(11) comment('权重')" json:"wgiht"`
	IfDel       int64  `xorm:"not null default 0 INT(11) comment('是否删除')" json:"if_del"`
}

//Category List
func ListCategorys(filter mysql.OrmFilter) (list []*Category, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

//Category Create
func CreateCategorys(message []*Category) (err error) {
	_, err = mysql.GetDB().Insert(message)
	return
}

//Category Update
func UpdateCategorys(filter mysql.OrmFilter, message []*Category) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Update(message)
	return
}
