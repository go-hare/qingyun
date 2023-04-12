package models

import (
	"qingyun/common/store/mysql"
)

type Env struct {
	mysql.Model  `xorm:"extends"`
	Env       string    `xorm:"not null default '' varchar(64)" json:"project"`
	Name      string    `xorm:"not null default '' varchar(64)" json:"alias"`
	Cluster   string    `xorm:"not null default '' varchar(64)" json:"cluster"`
	Addresses string    `xorm:"not null default '' varchar(128)" json:"addresses"`
	IfDefault int64     `xorm:"not null default 0 int" json:"ifDefault"`
}

func GetEnvList(filter mysql.OrmFilter) (list []Env, err error) {
	session := mysql.GetDB().NewSession()
	if filter != nil {
		session = filter(session)
	}

	err = session.Find(&list)
	return
}
