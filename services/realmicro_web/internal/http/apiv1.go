package http

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"qingyun/services/realmicro_web/internal/config"
	"qingyun/services/realmicro_web/internal/envs"
	models "qingyun/services/realmicro_web/models"
	"xorm.io/xorm"
)

const (
	Apiv1CodeOk = iota
	Apiv1CodeParamError
	Apiv1CodeInternalError
	Apiv1CodeLoginError
)

type Apiv1Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

type ConfigReq struct {
	Env     string `json:"env"`
	EnvName string `json:"envName"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	IfDir   bool   `json:"ifDir"`
}

func Apiv1Call(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	log.Infof("serve static file vars: %v", vars)
	method, ok := vars["method"]
	if !ok {
		log.Infof("cannot found method from request.")
		return
	}

	d := json.NewDecoder(r.Body)
	d.UseNumber()

	var rsp Apiv1Response
	defer func() {
		WriteJSON(w, http.StatusOK, rsp)
	}()
	switch method {
	case "edit.project":
		var req models.Project
		if err := d.Decode(&req); err != nil {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if req.Project == "" || req.Alias == "" {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if affected, err := models.UpdateProject(func(session *xorm.Session) *xorm.Session {
			return session.Where("project = ?", req.Project)
		}, &req); err != nil {
			log.Infof("update project error: %v", err)
			rsp.Code = Apiv1CodeInternalError
			return
		} else {
			if affected == 0 {
				if err = models.CreateProject(&req); err != nil {
					log.Infof("create project error: %v", err)
					rsp.Code = Apiv1CodeInternalError
					return
				}
			}
		}
	case "edit.app":
		var req models.App
		if err := d.Decode(&req); err != nil {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if req.AppName == "" || req.DepartmentId == 0 || req.Owner == 0 {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if affected, err := models.UpdateApp(func(session *xorm.Session) *xorm.Session {
			return session.Where("app_name = ?", req.AppName).Cols("desc", "department_id", "owner")
		}, &req); err != nil {
			log.Infof("update app error: %v", err)
			rsp.Code = Apiv1CodeInternalError
			return
		} else {
			if affected == 0 {
				if err = models.CreateApp(&req); err != nil {
					log.Infof("create app error: %v", err)
					rsp.Code = Apiv1CodeInternalError
					return
				}
			}
		}
	case "del.department":
		var req models.Department
		if err := d.Decode(&req); err != nil {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if req.Id == 0 {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if err := models.DelDepartment(func(session *xorm.Session) *xorm.Session {
			return session.ID(req.Id)
		}); err != nil {
			log.Infof("del department error: %v", err)
			rsp.Code = Apiv1CodeInternalError
			return
		}
	case "add.department":
		var req models.Department
		if err := d.Decode(&req); err != nil {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if req.Name == "" || req.NameEn == "" {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if err := models.CreateDepartment(&req); err != nil {
			log.Infof("add app error: %v", err)
			rsp.Code = Apiv1CodeInternalError
			return
		}
		rsp.Data = req
	case "del.admin":
		var req models.AdminUser
		if err := d.Decode(&req); err != nil {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if req.Id == 0 {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if err := models.DelAdminUser(func(session *xorm.Session) *xorm.Session {
			return session.ID(req.Id)
		}); err != nil {
			log.Infof("del admin user error: %v", err)
			rsp.Code = Apiv1CodeInternalError
			return
		}
	case "add.admin":
		var req models.AdminUser
		if err := d.Decode(&req); err != nil {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if req.Name == "" || req.Password == "" {
			rsp.Code = Apiv1CodeParamError
			return
		}
		req.Password = fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))
		if err := models.CreateAdminUser(&req); err != nil {
			log.Infof("add admin user error: %v", err)
			rsp.Code = Apiv1CodeInternalError
			return
		}
		rsp.Data = req
	case "del.config":
		var req ConfigReq
		if err := d.Decode(&req); err != nil {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if req.Key == "" {
			rsp.Code = Apiv1CodeParamError
			return
		}
		//if err := EtcdConfig.Del(req.Key); err != nil {
		if err := Envs.ConfigDel(envs.BuildEnv(req.Env, req.EnvName), req.Key); err != nil {
			log.Infof("del config error: %v", err)
			rsp.Code = Apiv1CodeInternalError
			return
		}
	case "put.config":
		var req ConfigReq
		if err := d.Decode(&req); err != nil {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if req.Key == "" {
			rsp.Code = Apiv1CodeParamError
			return
		}
		if req.IfDir {
			req.Value = config.DirValue
		}
		//if err := EtcdConfig.Put(req.Key, req.Value); err != nil {
		if err := Envs.ConfigPut(envs.BuildEnv(req.Env, req.EnvName), req.Key, req.Value); err != nil {
			log.Infof("put config error: %v", err)
			rsp.Code = Apiv1CodeInternalError
			return
		}
	default:
		log.Infof("cannot found this method: %s", method)
		return
	}
}
