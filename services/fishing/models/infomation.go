package models

import "qingyun/common/store/mysql"

type Infomation struct {
	mysql.Model          `xorm:"extends"`
	InfomationStatus     int64    `xorm:"not null default 0 INT(11) comment('状态 0，待审核 1，通过 -1 驳回')" json:"infomation_status"`
	InfomationType       int64    `xorm:"not null default 0 INT(11) comment('类型 0 钓点，1饵料')" json:"infomation_type"`
	AuthorTime           int64    `xorm:"not null default 0 INT(11) comment('审核时间')" json:"author_time"`
	RejectTime           int64    `xorm:"not null default 0 INT(11) comment('驳回时间 ')" json:"reject_time"`
	AuthorAdminId        int64    `xorm:"not null default 0 INT(11) comment('创建时间 ')" json:"author_admin_id"`
	Latitude             float64  `xorm:"not null default 0 DOUBLE comment('经度')" json:"latitude"`
	Longitude            float64  `xorm:"not null default 0 DOUBLE comment('精度')" json:"longitude"`
	UserId               int64    `xorm:"not null default 0 INT(11) comment('用户id') index(user_province_district_city)" json:"user_id"`
	Country              string   `xorm:"not null default '' VARCHAR(20) comment('国家') index(user_province_district_city)" json:"country"`
	City                 string   `xorm:"not null default '' VARCHAR(20) comment('市') index(user_province_district_city)" json:"city"`
	Province             string   `xorm:"not null default '' VARCHAR(20) comment('省') index(user_province_district_city)" json:"province"`
	District             string   `xorm:"not null default '' VARCHAR(20) comment('区') index(user_province_district_city)" json:"district"`
	Content              string   `xorm:"not null  TEXT  comment('内容')" json:"content"`
	Park                 string   `xorm:"not null default '' VARCHAR(256) comment('公园')" json:"park"`
	CreateTime           int64    `xorm:"not null default 0 INT(11) comment('创建时间 ')" json:"create_time"`
	Harvest              int64    `xorm:"not null default 0 INT(11) comment('收获 斤')" json:"harvest"`
	PayType              int64    `xorm:"not null default 0 INT(11) comment('支付类型 0 免费 1收费')" json:"pay_type"`
	Price                int64    `xorm:"not null default 0 INT(11) comment('金额')" json:"price"`
	ListTags             []string `xorm:"not null default '' VARCHAR(256) comment('标签')" json:"list_tags"`
	ListInfomationImages []string `xorm:"not null TEXT  json comment('鱼获图片')" json:"list_infomation_images"`
	ListEnvImages        []string `xorm:"not null  TEXT  json comment('环境图片')" json:"list_env_images"`
	ListBaits            []*Bait  `xorm:"not null  TEXT  json comment('环境图片')" json:"list_baits"`
	UpdateTime           int64    `xorm:"not null default 0 INT(11) comment('更新时间 ')" json:"update_time"`
	LikeCount            int64    `xorm:"not null default 0 INT(11) comment('点赞数量 ')" json:"like_count"`
	CommentCount         int64    `xorm:"not null default 0 INT(11) comment('评论')" json:"comment_count"`
	CollectCount         int64    `xorm:"not null default 0 INT(11) comment('收藏数量')" json:"collect_count"`
	ReadCount            int64    `xorm:"not null default 0 INT(11) comment('热度')" json:"read_count"`
	Hot                  int64    `xorm:"not null default 0 INT(11) comment('')" json:"hot"`
	IsDel                int64    `xorm:"not null default 0 INT(11) comment('是否删除 0否 1是')" json:"is_del"`
	DelTime              int64    `xorm:"not null default 0 INT(11) comment('是否删除')" json:"del_time"`
}

type Bait struct {
	Image       string `json:"image" `
	BaitName    string `json:"bait_name"`
	BaitPercent int64  `json:"bait_percent"`
}

func ListInfomations(filter mysql.OrmFilter) (list []*Infomation, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func GetInfomation(filter mysql.OrmFilter) (*Infomation, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(Infomation)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func UpdateInfomation(filter mysql.OrmFilter, message *Infomation) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Update(message)
	return
}
