package models

import "qingyun/common/store/mysql"

type ArticleTag struct {
	mysql.Model `xorm:"extends"`
	Name        string `xorm:"not null default '' VARCHAR(256) comment('名称')" json:"name"`
	Weight      int64  `xorm:"not null default 0 INT(11) comment('权重')" json:"weight"`
	Hot         int64  `xorm:"not null default 0 INT(11) comment('热度')" json:"hot"`
	CreateTime  int64  `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
	UsedCount   int64  `xorm:"not null default 0 INT(11) comment('参与')" json:"used_count"`
	ViewCount   int64  `xorm:"not null default 0 INT(11) comment('展示')" json:"view_count"`
	IsDel       int64  `xorm:"not null default 0 INT(11) comment('是否删除 0 否 1是')" json:"is_del"`
	DelTime     int64  `xorm:"not null default 0 INT(11) comment('删除时间')" json:"del_time"`
}

func GetArticleTag(filter mysql.OrmFilter) (*ArticleTag, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(ArticleTag)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func ListArticleTags(filter mysql.OrmFilter) (list []*ArticleTag, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func CreateArticleTag(message *ArticleTag) (err error) {
	_, err = mysql.GetDB().Insert(message)
	return
}

func UpdateArticleTag(filter mysql.OrmFilter, message *ArticleTag) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Update(message)
	return
}
