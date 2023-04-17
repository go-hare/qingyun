package models

import "qingyun/common/store/mysql"

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type Tutorial struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type Model struct {
	mysql.Model           `xorm:"extends"`
	CategoryId            int64                    `xorm:"not null default 0 INT(11) comment('分类id')" json:"category_id"`
	Name                  string                   `xorm:"not null default '' VARCHAR(256) comment('名称')" json:"name"`
	Icon                  string                   `xorm:"not null default '' VARCHAR(256) comment('图标')" json:"icon"`
	CreateTime            int64                    `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
	ChatCompletionMessage []*ChatCompletionMessage `xorm:"not null TEXT json comment('模型上下文')" json:"chat_completion_message"`
	Title                 string                   `xorm:"not null default '' VARCHAR(256) comment('标题')" json:"title"`
	Desc                  string                   `xorm:"not null default '' VARCHAR(256) comment('描述')" json:"desc"`
	ClickNumber           int64                    `xorm:"not null default 0 INT(11) comment('点击量')" json:"click_number"`
	Weight                int64                    `xorm:"not null default 0 INT(11) comment('权重')" json:"weight"`
	Tutorial              []*Tutorial              `xorm:"not null TEXT json comment('教程')" json:"tutorial"`
	IfDel                 int64                    `xorm:"not null default 0 INT(11) comment('是否删除')" json:"if_del"`
}

func ListModels(filter mysql.OrmFilter) (list []*Model, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func GetModel(filter mysql.OrmFilter) (user *Model, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	user = &Model{}
	_, err = session.Get(user)
	return
}
func CreateModels(message []*Model) (err error) {
	_, err = mysql.GetDB().Insert(message)
	return
}
