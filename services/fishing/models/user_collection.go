package models

import (
	"qingyun/common/store/mysql"
	"strconv"
)

const UserCollectionNum = 10

type UserCollection struct {
	mysql.Model        `xorm:"extends"`
	UserId             int64 `xorm:"not null default 0 INT(11) comment('用户id') unique(user_collection)" json:"user_id"`
	CollectionId       int64 `xorm:"not null default 0 INT(11) comment('收藏id') unique(user_collection)" json:"collection_id"`
	UserCollectionType int64 `xorm:"not null default 0 INT(11) comment('类型 0 信息 1技巧') unique(user_collection)" json:"user_collection_type"`
	CollectionTime     int64 `xorm:"not null default 0 INT(11) comment('收藏时间')" json:"collection_time"`
}

func (rc *UserCollection) TableName() string {
	return "user_collection_" + strconv.Itoa(int(rc.UserId)%UserCollectionNum)
}

func ListUserCollections(filter mysql.OrmFilter) (list []*UserCollection, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func GetUserCollection(filter mysql.OrmFilter) (*UserCollection, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(UserCollection)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func CreateUserCollection(message *UserCollection) (err error) {
	_, err = mysql.GetDB().Table(message.TableName()).Insert(message)
	return
}

func DeleteUserCollection(filter mysql.OrmFilter) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Delete(&UserCollection{})
	return
}
