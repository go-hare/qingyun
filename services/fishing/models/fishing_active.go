package models

import "qingyun/common/store/mysql"

type FishingActive struct {
	mysql.Model `xorm:"extends"`
	UserId      int64    `xorm:"not null default 0 INT(11) comment('用户id')" json:"user_id"`
	FreeType    int64    `xorm:"not null default 0 INT(11) comment('收费类型 0 免费 1收费')" json:"free_type"`
	FreePrice   int64    `xorm:"not null default 0 INT(11) comment('收费金额')" json:"free_price"`
	Describe    string   `xorm:"not null default '' VARCHAR(256) comment('描述')" json:"describe"`
	ImageList   []string `xorm:"not null TEXT json comment('图片列表')" json:"image_list"`
}
