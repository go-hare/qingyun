package http

import "net/http"

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	renderView(w, r, "templates/login.tpl", nil)
}
