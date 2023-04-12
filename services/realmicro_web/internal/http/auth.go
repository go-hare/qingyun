package http

import (
	"crypto/md5"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"
	"net/http"
	models "qingyun/services/realmicro_web/models"
	"xorm.io/xorm"
)

var cookieHandler = securecookie.New(
	[]byte("2HRDy6DOTFeO2STDHV2gk3BvzHjCmxHv"),
	[]byte("wKgsCLDxaIBFr9Jy"))

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("micro-web-session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("micro-web-session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("micro-web-session", value); err == nil {
		cookie := &http.Cookie{
			Name:   "micro-web-session",
			Value:  encoded,
			Path:   "/",
			MaxAge: 86400,
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "micro-web-session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func init() {
	//store = sessions.NewCookieStore(
	//	[]byte("micro-web-session"),
	//)
	gob.Register(&models.AdminUser{})

	//store.Options = &sessions.Options{
	//	MaxAge:   3600 * 24, // 24 hour
	//	HttpOnly: true,
	//}
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" || r.URL.Path == "/logout" || r.URL.Path == "/api/login" {
			next.ServeHTTP(w, r)
			return
		}
		userName := getUserName(r)
		if userName != "" {
			log.Infof("[auth middleware] url[%s] get user name : %s", r.URL.Path, userName)
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", 302)
		}
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.UseNumber()

	var rsp Apiv1Response
	defer func() {
		WriteJSON(w, http.StatusOK, rsp)
	}()

	var req models.AdminUser
	if err := d.Decode(&req); err != nil {
		rsp.Code = Apiv1CodeParamError
		return
	}
	if req.Name == "" || req.Password == "" {
		rsp.Code = Apiv1CodeParamError
		return
	}

	admin, err := models.GetAdminUser(func(session *xorm.Session) *xorm.Session {
		return session.Where("name = ?", req.Name)
	})
	if err != nil {
		log.Infof("get admin user error: %v", err)
		rsp.Code = Apiv1CodeInternalError
		return
	}
	if admin == nil {
		log.Infof("get admin of %s is nil", req.Name)
		rsp.Code = Apiv1CodeLoginError
		return
	}
	if fmt.Sprintf("%x", md5.Sum([]byte(req.Password))) != admin.Password {
		log.Infof("get admin of %s --- %s : password error", req.Name, req.Password)
		rsp.Code = Apiv1CodeLoginError
		return
	}
	setSession(admin.Name, w)
	log.Infof("login success with name: %s", admin.Name)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}
