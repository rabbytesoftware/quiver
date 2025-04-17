package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	packages_routes "rounds.com.ar/watcher/rest/service/packages"
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
	subrouter := router.PathPrefix("api-url/v1").Subrouter()
	packagesHandler := packages_routes.NewHandler().PackagesRoutes(subrouter)
	singlePackageHandler := packages_routes.NewHandler().SinglePackageRoutes(subrouter)

	log.Println("Server listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}