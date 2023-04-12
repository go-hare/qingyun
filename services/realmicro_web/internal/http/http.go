package http

import (
	"fmt"
	"github.com/realmicro/realmicro/registry"
	"github.com/serenize/snaker"
	"html/template"
	"net/http"
	"path/filepath"
	"qingyun/services/realmicro_web/internal/config"
	"qingyun/services/realmicro_web/internal/envs"
	"strings"
)

var (
	Envs *envs.PaaSEnvs
)

func InitEnvs() {
	Envs = envs.NewPaaSEnvs()
}

func format(v *registry.Value) string {
	if v == nil || len(v.Values) == 0 {
		return "{}"
	}
	var f []string
	for _, k := range v.Values {
		f = append(f, formatEndpoint(k, 0))
	}
	return fmt.Sprintf("{\n%s}", strings.Join(f, ""))
}

func formatEndpoint(v *registry.Value, r int) string {
	// default format is tabbed plus the value plus new line
	fparts := []string{"", "%s %s", "\n"}
	for i := 0; i < r+1; i++ {
		fparts[0] += "\t"
	}
	// its just a primitive of sorts so return
	if len(v.Values) == 0 {
		return fmt.Sprintf(strings.Join(fparts, ""), snaker.CamelToSnake(v.Name), v.Type)
	}

	// this thing has more things, it's complex
	fparts[1] += " {"

	vals := []interface{}{snaker.CamelToSnake(v.Name), v.Type}

	for _, val := range v.Values {
		fparts = append(fparts, "%s")
		vals = append(vals, formatEndpoint(val, r+1))
	}

	// at the end
	l := len(fparts) - 1
	for i := 0; i < r+1; i++ {
		fparts[l] += "\t"
	}
	fparts = append(fparts, "}\n")

	return fmt.Sprintf(strings.Join(fparts, ""), vals...)
}

func formatConfig(data map[string]*config.Node) template.HTML {
	if data == nil || len(data) == 0 {
		return ""
	}
	var results template.HTML
	for k, v := range data {
		if !v.IsDir {
			continue
		}
		results += template.HTML(fmt.Sprintf(`<ul id="ulNode-%s"><li class="ng-scope" id="liNode-%s"><div class="node" id="dirNode-%s" onclick="configValues('%s')">%s</div>
			<div class="btn-group ng-scope btn-group-hide" id="dirNodeBtn-%s">
            	<button class="btn btn-default btn-xs ng-scope btn-list" type="button" title="Create Directory" data-toggle="modal" data-target="#addConfigDirModal">
                	<span class="glyphicon glyphicon-plus"></span>
            	</button>
                <button type="button" class="btn btn-default btn-xs ng-scope btn-list" title="Delete Directory" onclick="delConfigDir('%s')">
                	<span class="glyphicon glyphicon-trash"></span>
                </button>
			</div>
		`, v.LongKey, v.LongKey, v.LongKey, v.LongKey, k, v.LongKey, v.LongKey))
		if len(v.Nodes) != 0 {
			results += formatConfig(v.Nodes)
		}
		results += `</li></ul>`
	}
	return template.HTML(results)
}

func getEnvs(env, envName string) (es []*envs.EnvDetail, defaultEnv *envs.Env) {
	es, defaultEnv = Envs.GetEnvs()
	if len(env) > 0 && len(envName) > 0 {
		for i := 0; i < len(es); i++ {
			for j := 0; j < len(es[i].Envs); j++ {
				if es[i].Envs[j].Env == env && es[i].Envs[j].Name == envName {
					defaultEnv = es[i].Envs[j]
					es[i].Envs[j].IfSelected = true
					break
				}
			}
		}
	} else {
		for i := 0; i < len(es); i++ {
			for j := 0; j < len(es[i].Envs); j++ {
				if es[i].Envs[j].Env == defaultEnv.Env && es[i].Envs[j].Name == defaultEnv.Name {
					es[i].Envs[j].IfSelected = true
					break
				}
			}
		}
	}
	return
}

func render(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	t, err := template.New("template").Funcs(template.FuncMap{
		"format": format,
		"fc":     formatConfig,
	}).ParseFiles("templates/layout.tpl", "templates/menu-env.tpl", "templates/menu-opr.tpl", tmpl)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	if err := t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Results":   data,
		"LoginUser": getUserName(r),
	}); err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
	}
}

func renderView(w http.ResponseWriter, r *http.Request, tpl string, data interface{}) {
	t := template.New(filepath.Base(tpl))
	//t = t.Funcs(template.FuncMap{"unescaped": unescaped, "formattime": formattime})
	t, err := t.ParseFiles(tpl)
	if err != nil {
		return
	}
	err = t.ExecuteTemplate(w, filepath.Base(tpl), data)
	if err != nil {
		return
	}
}
