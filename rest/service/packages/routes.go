package packages_routes

import (
	"github.com/gorilla/mux"
	logger "github.com/rabbytesoftware/quiver/logger"
	packages "github.com/rabbytesoftware/quiver/packages"
	mph "github.com/rabbytesoftware/quiver/rest/service/packages/handlers/multiple-packages"
	sph "github.com/rabbytesoftware/quiver/rest/service/packages/handlers/single-package"
)

type Handler struct {
	logs	*logger.Logger
	pkgs	*map[string]*packages.Package
}

func NewHandler(
	logs *logger.Logger,
	pkgs *map[string]*packages.Package,
) *Handler{
	return &Handler{
		logs: logs,
		pkgs: pkgs,
	}
}

func (h *Handler) PackagesRoutes(router *mux.Router) {
	mph := mph.NewPackagesHandler(h.logs, h.pkgs)

	// Get packages list 
	router.HandleFunc("/package", mph.GetPackagesList).Methods("GET")
}

func (h *Handler) SinglePackageRoutes(router *mux.Router) {
	sph := sph.NewPackageHandler(h.logs, h.pkgs)

	// Get single package data 
	router.HandleFunc("/package/{name}", sph.GetSinglePackageHandler).Methods("GET")

	// Delete package 
	router.HandleFunc("/package/{name}", sph.DeletePackageHandler).Methods("DELETE")

	// Install specific package 
	router.HandleFunc("/package/{name}", sph.InstallPackageHandler).Methods("POST")

	// Run package 
	router.HandleFunc("/package/{name}/run", sph.RunPackageHandler).Methods("PATCH")

	// Initialize package
	router.HandleFunc("/package/{name}/init", sph.InitPackageHandler)

	// Stop package process 
	router.HandleFunc("/package/{name}/stop", sph.ShutdownPackageHandler).Methods("PATCH")

}

