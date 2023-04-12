package models

import (
	"qingyun/common/store/mysql"
	"strconv"
)

const UserLikeNum = 10

type UserLike struct {
	mysql.Model  `xorm:"extends"`
	UserId       int64 `xorm:"not null default 0 INT(11) comment('用户id') unique(user_like)" json:"user_id"`
	LikeId       int64 `xorm:"not null default 0 INT(11) comment('点赞id') unique(user_like)" json:"like_id"`
	UserLikeType int64 `xorm:"not null default 0 INT(11) comment('类型 0 信息 1技巧') unique(user_like)" json:"user_like_type"`
	LikeTime     int64 `xorm:"not null default 0 INT(11) comment('点赞时间')" json:"like_time"`
}

func (rc *UserLike) TableName() string {
	return "user_like_" + strconv.Itoa(int(rc.UserId)%UserLikeNum)
}

func ListUserLikes(filter mysql.OrmFilter) (list []*UserLike, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func GetUserLike(filter mysql.OrmFilter) (*UserLike, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(UserLike)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func CreateUserLike(message *UserLike) (err error) {
	_, err = mysql.GetDB().Table(message.TableName()).Insert(message)
	return
}

func DeleteUserLike(filter mysql.OrmFilter) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Delete(&UserLike{})
	return
}
