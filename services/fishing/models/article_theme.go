package models

import "qingyun/common/store/mysql"

type ArticleTheme struct {
	mysql.Model `xorm:"extends"`
	CategoryId  int64  `xorm:"not null default 0 INT(11) comment('分类id')" json:"category_id"`
	TagId       int64  `xorm:"not null default 0 INT(11) comment('标签id')" json:"tag_id"`
	Title       string `xorm:"not null default '' VARCHAR(256) comment('标题')" json:"title"`
	Desc        string `xorm:"not null default '' VARCHAR(256) comment('描述')" json:"desc"`
	BannerlUrl  string `xorm:"not null default '' VARCHAR(256) comment('banner')" json:"bannerl_url"`
	ImgUrl      string `xorm:"not null default '' VARCHAR(256) comment('banner')" json:"img_url"`
	Weight      int64  `xorm:"not null default 0 INT(11) comment('权重')" json:"weight"`
	CreateTime  int64  `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
	IsDel       int64  `xorm:"not null default 0 INT(11) comment('是否删除 0 否 1是')" json:"is_del"`
	UsedCount   int64  `xorm:"not null default 0 INT(11) comment('参与')" json:"used_count"`
	ViewCount   int64  `xorm:"not null default 0 INT(11) comment('展示')" json:"view_count"`
	DelTime     int64  `xorm:"not null default 0 INT(11) comment('删除时间')" json:"del_time"`
}

func GetArticleTheme(filter mysql.OrmFilter) (*ArticleTheme, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(ArticleTheme)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func ListArticleThemes(filter mysql.OrmFilter) (list []*ArticleTheme, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func CreateArticleTheme(message *ArticleTheme) (err error) {
	_, err = mysql.GetDB().Insert(message)
	return
}
