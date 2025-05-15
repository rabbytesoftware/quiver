package api

import (
	"net/http"

	"github.com/gorilla/mux"
	packages_routes "rounds.com.ar/watcher/rest/service/packages"
	logger "rounds.com.ar/watcher/view/logger"
)

type ApiServer struct{
	addr string;
}

func CreateServerAPI(addr string) (*ApiServer){
	return &ApiServer{
			addr: addr,
		}
}

func (s *ApiServer) Run() error{
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	packages_routes.NewHandler().PackagesRoutes(subrouter)
	packages_routes.NewHandler().SinglePackageRoutes(subrouter)

	logger.It.Info("server-api", "Server listening on port %s", s.addr)

	return http.ListenAndServe(s.addr, router)
}