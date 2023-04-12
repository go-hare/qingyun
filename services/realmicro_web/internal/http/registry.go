package http

import (
	"encoding/json"
	"github.com/realmicro/realmicro/registry"
	"net/http"
	"qingyun/services/realmicro_web/internal/permission"
	models "qingyun/services/realmicro_web/models"
	"sort"
	"xorm.io/xorm"
)

func RegistryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	svc := r.Form.Get("service")
	env := r.Form.Get("env")
	envName := r.Form.Get("envname")

	es, defaultEnv := getEnvs(env, envName)

	if len(svc) > 0 {
		s, err := Envs.GetService(defaultEnv.Build(), svc)
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}

		if len(s) == 0 {
			http.Error(w, "Not found", 404)
			return
		}

		app, err := models.GetAppDetail(func(session *xorm.Session) *xorm.Session {
			return session.Where("app_name = ?", s[0].Name)
		})
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}
		if app == nil {
			app = &models.AppDetail{
				App: models.App{
					AppName: s[0].Name,
				},
			}
		}

		if r.Header.Get("Content-Type") == "application/json" {
			b, err := json.Marshal(map[string]interface{}{
				"services": s,
			})
			if err != nil {
				http.Error(w, "Error occurred:"+err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		}

		//svcConfig, err := EtcdConfig.Get(svc)
		svcConfig, err := Envs.GetSvcConfig(defaultEnv.Build(), svc)
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}

		render(w, r, "templates/service.tpl", map[string]interface{}{
			"Svc":       svc,
			"Services":  s,
			"AppDetail": app,
			"Config":    svcConfig,
			"Env":       defaultEnv.Env,
			"EnvName":   defaultEnv.Name,
		})
		return
	}

	var services []*registry.Service
	var err error
	services, err = Envs.ListServices(defaultEnv.Build())
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	classifyServices := make(map[string][]*registry.Service)
	var projects []string
	for i := 0; i < len(services); i++ {
		s, _ := Envs.GetService(defaultEnv.Build(), services[i].Name)
		if len(s) > 0 {
			if len(s[0].Nodes) > 0 {
				project, ok := s[0].Nodes[0].Metadata["project"]
				if ok {
					if _, ok := classifyServices[project]; !ok {
						projects = append(projects, project)
					}
					classifyServices[project] = append(classifyServices[project], s[0])
				} else {
					classifyServices["Others"] = append(classifyServices["Others"], s[0])
				}
			}
		}
	}

	sort.Strings(projects)
	projects = append(projects, "Others")
	//sort.Sort(sortedServices{services})
	projectList, err := models.GetProjectList(func(session *xorm.Session) *xorm.Session {
		return session.In("project", projects)
	})
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	type ProjectService struct {
		Project      string
		ProjectAlias string
		Services     []*registry.Service
	}
	projectServices := make([]ProjectService, len(projects))
	for i := 0; i < len(projects); i++ {
		p := projects[i]
		for j := 0; j < len(projectList); j++ {
			if projects[i] == projectList[j].Project {
				p = projectList[j].Alias + "(" + projectList[j].Project + ")"
				break
			}
		}
		projectServices[i] = ProjectService{
			Project:      projects[i],
			ProjectAlias: p,
			Services:     classifyServices[projects[i]],
		}
	}

	if r.Header.Get("Content-Type") == "application/json" {
		b, err := json.Marshal(map[string]interface{}{
			"Projects": projectServices,
		})
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}

	render(w, r, "templates/register.tpl", map[string]interface{}{
		"Projects": projectServices,
		"Envs":     es,
		"Uri":      "registry",
		"Env":      defaultEnv.Env,
		"EnvName":  defaultEnv.Name,
		"Menus":    permission.GetMenus(getUserName(r)),
	})
}
