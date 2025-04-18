package packages_routes

import (
	"github.com/gorilla/mux"
	multiple_packages_handlers "rounds.com.ar/watcher/rest/service/packages/handlers/multiple-packages"
	single_package_handlers "rounds.com.ar/watcher/rest/service/packages/handlers/single-package"
)

type Handler struct{

}

func NewHandler() *Handler{
	return &Handler{}
}

func (h *Handler) PackagesRoutes(router *mux.Router){
	// Get packages list 
	router.HandleFunc("/package", multiple_packages_handlers.GetPackagesList).Methods("GET")
}

func (h *Handler) SinglePackageRoutes(router *mux.Router){
	// Get single package data 
	router.HandleFunc("/package/{name}", single_package_handlers.GetSinglePackageHandler).Methods("GET")

	// Delete package 
	router.HandleFunc("/package/{name}", single_package_handlers.DeletePackageHandler).Methods("DELETE")

	// Install specific package 
	router.HandleFunc("/package/{name}", single_package_handlers.InstallPackageHandler).Methods("POST")

	// Run package 
	router.HandleFunc("/package/{name}/run", single_package_handlers.RunPackageHandler).Methods("PATCH")

	// Initialize package
	router.HandleFunc("/package/{name}/init", single_package_handlers.InitPackageHandler)

	// Stop package process 
	router.HandleFunc("/package/{name}/stop", single_package_handlers.ShutdownPackageHandler).Methods("PATCH")

}

