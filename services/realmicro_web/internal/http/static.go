package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func ServeStaticFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Infof("serve static file vars: %v", vars)
	typename, ok := vars["typename"]
	if !ok {
		log.Infof("cannot found typename from request.")
		return
	}
	filename, ok := vars["filename"]
	if !ok {
		log.Infof("cannot found filename from request.")
		return
	}

	http.ServeFile(w, r, fmt.Sprintf("./static/%s/%s", typename, filename))
}

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/img/favicon.ico")
}
