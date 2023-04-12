package permission

import (
	log "github.com/sirupsen/logrus"
	"qingyun/services/realmicro_web/internal/roles"
	models "qingyun/services/realmicro_web/models"
	"xorm.io/xorm"
)

type Menu struct {
	Icon       string
	Name       string
	Url        string
	Permission string
}

var (
	Menus = []*Menu{
		&Menu{
			Icon:       "icon-test-manage",
			Name:       "Call Test",
			Url:        "/call",
			Permission: roles.RoleDev,
		},
		&Menu{
			Icon:       "icon-project-manage",
			Name:       "服务管理",
			Url:        "/app",
			Permission: roles.RoleAdmin,
		},
		&Menu{
			Icon:       "icon-bumen-manage",
			Name:       "部门管理",
			Url:        "/department",
			Permission: roles.RoleGod,
		},
		&Menu{
			Icon:       "icon-admin-manage",
			Name:       "人员管理",
			Url:        "/admin",
			Permission: roles.RoleGod,
		},
	}
)

func GetMenus(user string) (ms []*Menu) {
	admin, err := models.GetAdminUser(func(session *xorm.Session) *xorm.Session {
		return session.Where("name = ?", user)
	})
	if err != nil {
		log.Info("get admin user error: %v", err)
		return
	}
	if admin == nil {
		return
	}
	for i := 0; i < len(Menus); i++ {
		if CheckPermission(admin, Menus[i].Permission) {
			ms = append(ms, Menus[i])
		}
	}
	return
}
