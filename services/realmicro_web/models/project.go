package models

import (
	"qingyun/common/store/mysql"
)

type Project struct {
	mysql.Model  `xorm:"extends"`
	Project   string    `xorm:"not null default '' varchar(64) unique" json:"project"`
	Alias     string    `xorm:"not null default '' varchar(64)" json:"alias"`
}

func CreateProject(info *Project) (err error) {
	_, err = mysql.GetDB().Insert(info)
	return
}

func UpdateProject(filter mysql.OrmFilter, info *Project) (affected int64, err error) {
	session := mysql.GetDB().NewSession()
	if filter != nil {
		session = filter(session)
	}

	if info == nil {
		info = new(Project)
	}
	affected, err = session.Update(info)
	return
}

func GetProject(filter mysql.OrmFilter) (*Project, error) {
	session := mysql.GetDB().NewSession()
	if filter != nil {
		session = filter(session)
	}

	info := new(Project)
	if has, err := session.Get(info); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return info, nil
	}
}

func GetProjectList(filter mysql.OrmFilter) (list []Project, err error) {
	session := mysql.GetDB().NewSession()
	if filter != nil {
		session = filter(session)
	}

	err = session.Find(&list)
	return
}
