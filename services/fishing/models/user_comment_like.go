package models

import (
	"qingyun/common/store/mysql"
	"strconv"
)

const UserCommentLikeNum = 10

type UserCommentLike struct {
	mysql.Model `xorm:"extends"`
	UserId      int64 `xorm:"not null default 0 INT(11) comment('用户id') unique(user_like)" json:"user_id"`
	CommentId   int64 `xorm:"not null default 0 INT(11) comment('点赞id') unique(user_like)" json:"comment_id"`
	LikeTime    int64 `xorm:"not null default 0 INT(11) comment('点赞时间')" json:"like_time"`
}

func (rc *UserCommentLike) TableName() string {
	return "user_comment_like_" + strconv.Itoa(int(rc.UserId)%UserCommentLikeNum)
}

func CreateUserCommentLike(message *UserCommentLike) (err error) {
	_, err = mysql.GetDB().Table(message.TableName()).Insert(message)
	return
}

func ListUserCommentLikes(filter mysql.OrmFilter) (list []*UserCommentLike, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func GetUserCommentLike(filter mysql.OrmFilter) (*UserCommentLike, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(UserCommentLike)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func DeleteUserCommentLike(filter mysql.OrmFilter) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Delete(&UserCommentLike{})
	return
}
