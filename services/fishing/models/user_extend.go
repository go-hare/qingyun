package models

import "qingyun/common/store/mysql"

type UserExtend struct {
	mysql.Model     `xorm:"extends"`
	UserId          int64 `xorm:"not null default 0 INT(11) comment('用户id') unique(user_id)" json:"user_id"`
	FollowCount     int64 `xorm:"not null default 0 INT(11) comment('关注数量')" json:"follow_count"`
	FansCount       int64 `xorm:"not null default 0 INT(11) comment('被关注数量')" json:"fans_count"`
	InviteCount     int64 `xorm:"not null default 0 INT(11) comment('邀请')" json:"invite_count"`
	UnReadCount     int64 `xorm:"not null default 0 INT(11) comment('未读数量')" json:"un_read_count"`
	InfomationCount int64 `xorm:"not null default 0 INT(11) comment('鱼获数量')" json:"infomation_count"`
	ArticleCount    int64 `xorm:"not null default 0 INT(11) comment('文章数量')" json:"article_count"`
}

func GetUserExtend(filter mysql.OrmFilter) (*UserExtend, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(UserExtend)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func ListUserExtends(filter mysql.OrmFilter) (list []*UserExtend, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func UpdateUserExtend(filter mysql.OrmFilter, message *UserExtend) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Update(message)
	return
}
