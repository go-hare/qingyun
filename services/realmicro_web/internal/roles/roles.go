package roles

const (
	RoleGod   = "god"
	RoleAdmin = "admin"
	RoleDev   = "dev"
)

var (
	Roles    = []string{RoleGod, RoleAdmin, RoleDev}
	RoleDesc = map[string]string{
		RoleGod:   "超级管理员",
		RoleAdmin: "管理员",
		RoleDev:   "开发者",
	}
)
