package api

import (
	"net/http"

	"github.com/gorilla/mux"
	logger "github.com/rabbytesoftware/quiver/logger"
	packages "github.com/rabbytesoftware/quiver/packages"
	packagesRoutes "github.com/rabbytesoftware/quiver/rest/service/packages"
)

type ApiServer struct{
	addr string
	logs *logger.Logger
	pkgs *map[string]*packages.Package
}

func CreateServerAPI(
	addr string,
	pkgs *map[string]*packages.Package,
) (*ApiServer){
	return &ApiServer{
		addr: addr,
		logs: logger.NewLogger("API"),
		pkgs: pkgs,
	}
}

func (s *ApiServer) Run() error{
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	packagesRoutes.NewHandler(s.logs, s.pkgs).PackagesRoutes(subrouter)
	packagesRoutes.NewHandler(s.logs, s.pkgs).SinglePackageRoutes(subrouter)

	s.logs.Info("Server listening on port %s", s.addr)

	return http.ListenAndServe(s.addr, router)
}