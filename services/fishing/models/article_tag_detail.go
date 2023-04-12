package models

import "qingyun/common/store/mysql"

type ArticleTagDetail struct {
	mysql.Model  `xorm:"extends"`
	ArticleId    int64 `xorm:"not null default 0 INT(11) comment('文章id')" json:"article_id"`
	ArticleTagId int64 `xorm:"not null default 0 INT(11) comment('标签id')" json:"article_tag_id"`
	CreateTime   int64 `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
}

func ListArticleTagDetails(filter mysql.OrmFilter) (list []*ArticleTagDetail, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func CreateListArticleTagDetails(message []*ArticleTagDetail) (err error) {
	_, err = mysql.GetDB().Insert(message)
	return
}
