package permission

import (
	"qingyun/services/realmicro_web/internal/roles"
	models "qingyun/services/realmicro_web/models"
)

func CheckPermission(admin *models.AdminUser, role string) bool {
	for i := 0; i < len(roles.Roles); i++ {
		if admin.Role == roles.Roles[i] {
			return true
		}
		if roles.Roles[i] == role {
			return false
		}
	}
	return false
}
