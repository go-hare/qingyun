package models

import (
	"qingyun/common/store/mysql"
)

type Department struct {
	mysql.Model  `xorm:"extends"`
	Name      string    `xorm:"not null default '' varchar(64) unique" json:"name"`
	NameEn    string    `xorm:"not null default '' varchar(64)" json:"nameEn"`
}

func CreateDepartment(info *Department) (err error) {
	_, err = mysql.GetDB().Insert(info)
	return
}

func GetDepartmentList(filter mysql.OrmFilter) (list []Department, err error) {
	session := mysql.GetDB().NewSession()
	if filter != nil {
		session = filter(session)
	}

	err = session.Find(&list)
	return
}

func UpdateDepartment(filter mysql.OrmFilter, info *Department) (affected int64, err error) {
	session := mysql.GetDB().NewSession()
	if filter != nil {
		session = filter(session)
	}

	if info == nil {
		info = new(Department)
	}
	affected, err = session.Update(info)
	return
}

func DelDepartment(filter mysql.OrmFilter) (err error) {
	session := mysql.GetDB().NewSession()
	if filter != nil {
		session = filter(session)
	}

	_, err = session.Delete(&Department{})
	return
}

type AdminUser struct {
	mysql.Model  `xorm:"extends"`
	Name         string    `xorm:"not null default '' varchar(64) unique" json:"name"`
	NameCn       string    `xorm:"not null default '' varchar(64)" json:"nameCn"`
	Password     string    `xorm:"not null default '' varchar(64)" json:"password"`
	DepartmentId int64     `xorm:"not null default 0 int" json:"departmentId"`
	Role         string    `xorm:"not null default '' varchar(64)" json:"role"`
	Phone        string    `xorm:"not null default '' varchar(64)" json:"phone"`
	Wechat       string    `xorm:"not null default '' varchar(64)" json:"wechat"`
}

func CreateAdminUser(info *AdminUser) (err error) {
	_, err = mysql.GetDB().Insert(info)
	return
}

func GetAdminUser(filter mysql.OrmFilter) (*AdminUser, error) {
	session := mysql.GetDB().NewSession()
	if filter != nil {
		session = filter(session)
	}

	info := new(AdminUser)
	if has, err := session.Get(info); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return info, nil
	}
}

func GetAdminUserList(filter mysql.OrmFilter) (list []AdminUser, err error) {
	session := mysql.GetDB().NewSession()
	if filter != nil {
		session = filter(session)
	}

	err = session.Find(&list)
	return
}

func UpdateAdminUser(filter mysql.OrmFilter, info *AdminUser) (affected int64, err error) {
	session := mysql.GetDB().NewSession()
	if filter != nil {
		session = filter(session)
	}

	if info == nil {
		info = new(AdminUser)
	}
	affected, err = session.Update(info)
	return
}

func DelAdminUser(filter mysql.OrmFilter) (err error) {
	session := mysql.GetDB().NewSession()
	if filter != nil {
		session = filter(session)
	}

	_, err = session.Delete(&AdminUser{})
	return
}

type AdminUserDetail struct {
	AdminUser  `xorm:"extends"`
	Department `xorm:"extends"`
}

func (AdminUserDetail) TableName() string {
	return "admin_user"
}

func GetAdminUserDetailList(filter mysql.OrmFilter) (list []AdminUserDetail, err error) {
	session := mysql.GetDB().Join("LEFT", "department", "admin_user.department_id = department.id")
	if filter != nil {
		session = filter(session)
	}

	if err = session.Find(&list); err != nil {
		return
	}
	return
}
