package models

import (
	"qingyun/common/store/mysql"
	"strconv"
)

const MessageNum = 10

type Message struct {
	mysql.Model `xorm:"extends"`
	UserId      int64  `xorm:"not null default 0 INT(11) comment('用户id')" json:"user_id"`
	CreateTime  int64  `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
	ReadStatus  int64  `xorm:"not null default 0 INT(11) comment('类型 0 未读，1已读')" json:"read_status"`
	Title       string `xorm:"not null default '' VARCHAR(256) comment('标题')" json:"title"`
	Content     string `xorm:"not null default '' VARCHAR(256) comment('标题')" json:"content"`
	Avatar      string `xorm:"not null default '' VARCHAR(256) comment('发送者头像')" json:"avatar"`
}

func (rc *Message) TableName() string {
	return "message_" + strconv.Itoa(int(rc.UserId)%MessageNum)
}

func GetMessage(filter mysql.OrmFilter) (*Message, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(Message)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func ListMessages(filter mysql.OrmFilter) (list []*Message, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func UpdateMessage(filter mysql.OrmFilter, message *Message) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Update(message)
	return
}
