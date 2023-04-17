package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/realmicro/realmicro"
	mconfig "github.com/realmicro/realmicro/config"
	"github.com/realmicro/realmicro/config/encoder/yaml"
	"github.com/realmicro/realmicro/config/reader"
	"github.com/realmicro/realmicro/config/reader/json"
	"github.com/realmicro/realmicro/config/source/file"
	"github.com/realmicro/realmicro/logger"
	mlogrus "github.com/realmicro/realmicro/logger/logrus"
	"github.com/realmicro/realmicro/registry"
	"github.com/realmicro/realmicro/registry/etcd"
	"github.com/realmicro/realmicro/wrapper/select/dc"
	"github.com/sirupsen/logrus"
	"net/http"
	"qingyun/common/wrapper"
	"qingyun/services/realmicro_web/internal/config"
	mhttp "qingyun/services/realmicro_web/internal/http"
	"qingyun/services/realmicro_web/internal/logic"
	"qingyun/services/realmicro_web/internal/svc"
	micro_web "qingyun/services/realmicro_web/proto"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

var (
	CORS = map[string]bool{"*": true}
)

type Server struct {
	*mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); CORS[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else if len(origin) > 0 && CORS["*"] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		return
	}

	s.Router.ServeHTTP(w, r)
}

func main() {
	flag.Parse()
	c, _ := mconfig.NewConfig(
		mconfig.WithReader(
			json.NewReader(
				reader.WithEncoder(yaml.NewEncoder()),
			),
		),
	)
	var err error
	// load the config from a file source
	if err = c.Load(file.NewSource(
		file.WithPath(*configFile),
	)); err != nil {
		fmt.Println(err)
		return
	}
	var cfg config.Config
	if err = c.Scan(&cfg); err != nil {
		fmt.Println(err)
		return
	}
	logger.DefaultLogger = mlogrus.NewLogger(mlogrus.WithJSONFormatter(&logrus.JSONFormatter{}))

	var h http.Handler
	r := mux.NewRouter()
	s := &Server{r}
	h = s
	s.HandleFunc("/", mhttp.RegistryHandler)
	s.HandleFunc("/registry", mhttp.RegistryHandler)
	s.HandleFunc("/call", mhttp.CallHandler)
	s.HandleFunc("/department", mhttp.DepartmentHandler)
	s.HandleFunc("/admin", mhttp.AdminUserHandler)
	s.HandleFunc("/app", mhttp.AppHandler)
	s.HandleFunc("/app-edit", mhttp.AppEditHandler)
	s.HandleFunc("/rpc", mhttp.RPC)
	r.Handle("/api/v1/{method}", http.HandlerFunc(mhttp.Apiv1Call))
	r.Handle("/static/{typename}/{filename}", http.HandlerFunc(mhttp.ServeStaticFile))
	r.Handle("/favicon.ico", http.HandlerFunc(mhttp.ServeFavicon))

	s.HandleFunc("/login", mhttp.LoginPageHandler)
	s.HandleFunc("/api/login", mhttp.LoginHandler)
	s.HandleFunc("/logout", mhttp.LogoutHandler)

	r.Use(mhttp.Auth)

	srv := mhttp.NewServer(cfg.Address)
	srv.Init()
	srv.Handle("/", h)

	service := realmicro.NewService(
		realmicro.Name(cfg.ServerName),
		realmicro.Version(cfg.Version),
		realmicro.Registry(etcd.NewRegistry(registry.Addrs(cfg.Etcd.Address...))),
		realmicro.WrapHandler(wrapper.LogHandler()),
		realmicro.WrapClient(dc.NewDCWrapper, wrapper.LogCall),
	)
	service.Init()

	ctx := svc.NewServiceContext(&cfg, service)

	if err = micro_web.RegisterMicroWebHandler(service.Server(), logic.NewMicroWeb(ctx)); err != nil {
		logger.Fatal(err)
	}

	if err = srv.Start(); err != nil {
		logger.Fatal(err)
	}

	// Run server
	if err = service.Run(); err != nil {
		logger.Fatal(err)
	}

	if err = srv.Stop(); err != nil {
		logger.Fatal(err)
	}


}
