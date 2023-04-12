package http

import (
	"net/http"
	"qingyun/services/realmicro_web/internal/roles"

	log "github.com/sirupsen/logrus"
	models "qingyun/services/realmicro_web/models"
	"xorm.io/xorm"
)

func AppHandler(w http.ResponseWriter, r *http.Request) {
	appList, err := models.GetAppDetailList(nil)
	if err != nil {
		log.Infof("get app detail list error: %v", err)
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	myAppList, err := models.GetAppDetailList(func(session *xorm.Session) *xorm.Session {
		return session.Where("admin_user.name = ?", getUserName(r))
	})
	if err != nil {
		log.Infof("get my app detail list error: %v", err)
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	render(w, r, "templates/app.tpl", map[string]interface{}{
		"Apps":   appList,
		"MyApps": myAppList,
	})
}

func DepartmentHandler(w http.ResponseWriter, r *http.Request) {
	departments, err := models.GetDepartmentList(nil)
	if err != nil {
		log.Infof("get department list error: %v", err)
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	render(w, r, "templates/department.tpl", departments)
}

func AdminUserHandler(w http.ResponseWriter, r *http.Request) {
	admins, err := models.GetAdminUserDetailList(nil)
	if err != nil {
		log.Infof("get admin user list error: %v", err)
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}
	for i := 0; i < len(admins); i++ {
		if v, ok := roles.RoleDesc[admins[i].Role]; ok {
			admins[i].Role = v
		}
	}

	departments, err := models.GetDepartmentList(nil)
	if err != nil {
		log.Infof("get department list error: %v", err)
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	type Role struct {
		Role     string
		RoleDesc string
	}
	rs := make([]*Role, len(roles.Roles))
	for i := 0; i < len(roles.Roles); i++ {
		rs[i] = &Role{
			Role:     roles.Roles[i],
			RoleDesc: roles.RoleDesc[roles.Roles[i]],
		}
	}

	render(w, r, "templates/admin.tpl", map[string]interface{}{
		"Admins":      admins,
		"Departments": departments,
		"Roles":       rs,
	})
}

func AppEditHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	svc := r.Form.Get("service")
	if len(svc) == 0 {
		http.Error(w, "service cannot be nil", 500)
		return
	}

	departments, err := models.GetDepartmentList(nil)
	if err != nil {
		log.Infof("get department list error: %v", err)
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}
	admins, err := models.GetAdminUserList(nil)
	if err != nil {
		log.Infof("get admin list error: %v", err)
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	render(w, r, "templates/app-edit.tpl", map[string]interface{}{
		"Service":     svc,
		"Departments": departments,
		"Admins":      admins,
	})
}
