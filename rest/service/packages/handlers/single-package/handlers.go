package single_package_handlers

import (
	"net/http"
	"github.com/gorilla/mux"
)

func GetSinglePackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := params["name"]

	if len(pkgName) == 0 {
		http.Error(w, "Package name is required.", http.StatusBadRequest)
		return
	}
}

func DeletePackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := params["name"]

	if len(pkgName) == 0 {
		http.Error(w,"Package name is required.", http.StatusBadRequest)
		return
	}
}

func InstallPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := params["name"]

	if len(pkgName) == 0 {
		http.Error(w,"Package name is required.", http.StatusBadRequest)
		return
	}
}

func RunPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := params["name"]

	if len(pkgName) == 0 {
		http.Error(w, "Package name is required.", http.StatusBadRequest)
		return
	}
}

func ShutdownPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := params["name"]

	if len(pkgName) == 0 {
		http.Error(w, "Package name is required.", http.StatusBadRequest)
		return
	}
}