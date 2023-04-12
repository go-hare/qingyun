package models

import "qingyun/common/store/mysql"

type Banner struct {
	mysql.Model `xorm:"extends"`
	CreateTime  int64  `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
	BannerType  int64  `xorm:"not null default 0 INT(11) comment('类型 0 鱼获，1版块 ')" json:"banner_type"`
	ImageUrl    string `xorm:"not null default '' VARCHAR(256) comment('图片url')" json:"image_url"`
	ClickNumber int64  `xorm:"not null default 0 INT(11) comment('点击数量')" json:"click_number"`
	ViewType    int64  `xorm:"not null default 0 INT(11) comment('跳转类型 0, 鱼获 1，文章 ，2 主题 ,3，充值，4,webview，5 发布钓点 6，发布鱼饵，7 发布文章，8，我的账户')" json:"view_type"`
	ExternalId  int64  `xorm:"not null default 0 INT(11) comment('外部id,鱼获')" json:"external_id"`
	Title       string `xorm:"not null default '' VARCHAR(256) comment('标题')" json:"title"`
	LinkUrl     string `xorm:"not null default '' VARCHAR(256) comment('外部链接')" json:"link_url"`
	ThemeId     int64  `xorm:"not null default 0 INT(11) comment('主题id')" json:"theme_id"`
	Weight      int64  `xorm:"not null default 0 INT(11) comment('权重')" json:"weight"`
	IsDel       int64  `xorm:"not null default 0 INT(11) comment('是否删除 0 否 1是')" json:"is_del"`
	DelTime     int64  `xorm:"not null default 0 INT(11) comment('删除时间')" json:"del_time"`
}

func ListBanners(filter mysql.OrmFilter) (list []*Banner, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}
