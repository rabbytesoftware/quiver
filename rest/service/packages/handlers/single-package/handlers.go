package single_package_handlers

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	packages "rounds.com.ar/watcher/packages"
	pc "rounds.com.ar/watcher/sdk/base/package-config"
)

func getPackage(name string) *packages.Package {
	pkg := &packages.Package{
	PackageConfig: &pc.PackageConfig{
		// Inicializá campos necesarios acá
	},
	Runtime: &packages.PackageRuntime{
		// Inicializá campos necesarios acá
	},
}	

	return pkg
}

func checkIfPackageNameIsValid(name string, w http.ResponseWriter) {
	if len(name) == 0 {
		http.Error(w,"Package name is required.", http.StatusBadRequest)
		return
	}

	// Check if package name contains
	// whitespace.
	for _, char := range name {
		if char == ' ' {
			http.Error(w, "Invalid package name. Package names cannot contain spaces; use hyphens (-) instead.", http.StatusBadRequest)
			return
		}
	}

}

func GetSinglePackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	checkIfPackageNameIsValid(pkgName, w)

}

func DeletePackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	checkIfPackageNameIsValid(pkgName, w)
}

func InstallPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	checkIfPackageNameIsValid(pkgName, w)

	pkg := getPackage(pkgName)

	if _, err := pkg.Install(); err != nil {
		http.Error(w, "An error occurred while installing the package.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("The package has been successfully installed."))
}

func RunPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	checkIfPackageNameIsValid(pkgName, w)

	pkg := getPackage(pkgName)

	if _, err := pkg.Run(); err != nil {
		http.Error(w, "Error running package.", http.StatusBadRequest)
		return
	}
}

func ShutdownPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])


	checkIfPackageNameIsValid(pkgName, w)

	pkg := getPackage(pkgName)

	if err := pkg.Shutdown(); err != nil {
		http.Error(w, "Error shutting down package process.", http.StatusBadRequest)
		return
	}
}

func InitPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	checkIfPackageNameIsValid(pkgName, w)

	pkg := getPackage(pkgName)

	if err := pkg.Init(); err != nil {
		http.Error(w, "Error initializing package.", http.StatusBadRequest)
		return
	}
}