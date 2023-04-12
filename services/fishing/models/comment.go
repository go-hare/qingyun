package models

import "qingyun/common/store/mysql"

type Comment struct {
	mysql.Model      `xorm:"extends"`
	CommentType      int64  `xorm:"not null default 0 INT(11) comment('类型 0 信息 1文章')" json:"comment_type"`
	CommentStatus    int64  `xorm:"not null default 0 INT(11) comment('状态 0，待审核 1，通过 -1 驳回')" json:"comment_status"`
	CreateTime       int64  `xorm:"not null default 0 INT(11) comment('创建时间 ')" json:"create_time"`
	ObjectId         int64  `xorm:"not null default 0 INT(11) comment('对象id ')" json:"object_id"`
	ReleaseUserId    int64  `xorm:"not null default 0 INT(11) comment('发布用户 ')" json:"release_user_id"`
	LikeCount        int64  `xorm:"not null default 0 INT(11) comment('点赞数量')" json:"like_count"`
	RepliedCount     int64  `xorm:"not null default 0 INT(11) comment('回复数量')" json:"replied_count"`
	RepliedUserId    int64  `xorm:"not null default 0 INT(11) comment('回复用户id')" json:"replied_user_id"`
	RepliedCommentId int64  `xorm:"not null default 0 INT(11) comment('回复评论id')" json:"replied_comment_id"`
	Content          string `xorm:"not null default '' VARCHAR(256) comment('内容')" json:"content"`
}

func ListComments(filter mysql.OrmFilter) (list []*Comment, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func GetComment(filter mysql.OrmFilter) (*Comment, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(Comment)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func CreateComment(message *Comment) (err error) {
	_, err = mysql.GetDB().Insert(message)
	return
}

func UpdateComment(filter mysql.OrmFilter, message *Comment) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Update(message)
	return
}
