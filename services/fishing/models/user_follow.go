package models

import "qingyun/common/store/mysql"

type UserFollow struct {
	mysql.Model  `xorm:"extends"`
	UserId       int64 `xorm:"not null default 0 INT(11) comment('用户ID') unique(user_follow)" json:"user_id"`
	FollowUserId int64 `xorm:"not null default 0 INT(11) comment('关注用户') unique(user_follow)" json:"follow_user_id"`
	FollowTime   int64 `xorm:"not null default 0 INT(11) comment('关注时间')" json:"follow_time"`
}

func GetUserFollow(filter mysql.OrmFilter) (*UserFollow, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(UserFollow)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func CreateUserFollow(message *UserFollow) (err error) {
	_, err = mysql.GetDB().Insert(message)
	return
}

func ListUserFollows(filter mysql.OrmFilter) (list []*UserFollow, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func ListUserFollowsCount(filter mysql.OrmFilter) (list []*UserFollow, count int64, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}

	list = make([]*UserFollow, 0)
	count, err = session.FindAndCount(&list)
	return
}

func DeleteUserFollow(filter mysql.OrmFilter) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Delete(&UserFollow{})
	return
}
