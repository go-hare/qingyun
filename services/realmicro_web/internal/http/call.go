package http

import (
	"encoding/json"
	"github.com/realmicro/realmicro/registry"
	"net/http"
	"sort"
)

func CallHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	env := r.Form.Get("env")
	envName := r.Form.Get("envname")
	es, defaultEnv := getEnvs(env, envName)

	services, err := Envs.ListServices(defaultEnv.Build())
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	sort.Sort(sortedServices{services})

	serviceMap := make(map[string][]*registry.Endpoint)
	for _, service := range services {
		s, err := Envs.GetService(defaultEnv.Build(), service.Name)
		if err != nil {
			continue
		}
		if len(s) == 0 {
			continue
		}
		serviceMap[service.Name] = s[0].Endpoints
	}

	if r.Header.Get("Content-Type") == "application/json" {
		b, err := json.Marshal(map[string]interface{}{
			"services": services,
		})
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}

	render(w, r, "templates/call.tpl", map[string]interface{}{
		"ServiceMap": serviceMap,
		"Envs":       es,
		"Uri":        "call",
		"Env":        defaultEnv.Env,
		"EnvName":    defaultEnv.Name,
	})
}
