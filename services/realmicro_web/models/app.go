package models

import (
	"qingyun/common/store/mysql"
)

type App struct {
	mysql.Model  `xorm:"extends"`
	AppName      string    `xorm:"not null default '' varchar(64) unique" json:"appName"`
	Desc         string    `xorm:"not null default '' varchar(128)" json:"desc"`
	DepartmentId int64     `xorm:"not null default 0 int" json:"departmentId"`
	Owner        int64     `xorm:"not null default 0 int" json:"owner"`

}

func CreateApp(info *App) (err error) {
	_, err = mysql.GetDB().Insert(info)
	return
}

func UpdateApp(filter mysql.OrmFilter, info *App) (affected int64, err error) {
	session := mysql.GetDB().NewSession()
	if filter != nil {
		session = filter(session)
	}

	if info == nil {
		info = new(App)
	}
	affected, err = session.Update(info)
	return
}

func GetApp(filter mysql.OrmFilter) (*App, error) {
	session := mysql.GetDB().NewSession()
	if filter != nil {
		session = filter(session)
	}

	info := new(App)
	if has, err := session.Get(info); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return info, nil
	}
}

type AppDetail struct {
	App        `xorm:"extends"`
	AdminUser  `xorm:"extends"`
	Department `xorm:"extends"`
}

func (AppDetail) TableName() string {
	return "app"
}

func GetAppDetail(filter mysql.OrmFilter) (info *AppDetail, err error) {
	session := mysql.GetDB().Join("LEFT", "admin_user", "app.owner = admin_user.id").
		Join("LEFT", "department", "app.department_id = department.id")
	if filter != nil {
		session = filter(session)
	}

	info = new(AppDetail)
	if has, err := session.Get(info); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return info, nil
	}
}

func GetAppDetailList(filter mysql.OrmFilter) (list []AppDetail, err error) {
	session := mysql.GetDB().Join("LEFT", "admin_user", "app.owner = admin_user.id").
		Join("LEFT", "department", "app.department_id = department.id")
	if filter != nil {
		session = filter(session)
	}

	err = session.Find(&list)
	return
}
