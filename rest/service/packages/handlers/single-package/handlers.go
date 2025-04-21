package single_package_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	packages "rounds.com.ar/watcher/packages"
	packages_global_variables "rounds.com.ar/watcher/rest/shared/utils/packages/global-variables"
)

func getPackage(name string) (*packages.Package, error) {
	packagesList := packages_global_variables.Packages
	pkg, ok := packagesList[name]

	if !ok {
    return nil, fmt.Errorf("package '%s' not found", name)
	}

  return pkg, nil
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

	pkg, err := getPackage(pkgName)

	if err != nil {
		fmt.Println("Error:", err)
    return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pkg)
}

func DeletePackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	checkIfPackageNameIsValid(pkgName, w)

	// pkg, err := getPackage(pkgName)

	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// if _, err := pkg.Uninstall(); err != nil {
	// 	http.Error(w, "An error occurred while uninstalling the package.", http.StatusBadRequest)
	// 	return
	// }

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

func InstallPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	checkIfPackageNameIsValid(pkgName, w)

	pkg, err := getPackage(pkgName)

	if err != nil {
		fmt.Println("Error:", err)
    return
	}

	if _, err := pkg.Install(); err != nil {
		http.Error(w, "An error occurred while installing the package.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("The package has been successfully installed."))
}

func RunPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	checkIfPackageNameIsValid(pkgName, w)

	pkg, err := getPackage(pkgName)

	if err != nil {
		fmt.Println("Error:", err)
    return
	}

	if _, err := pkg.Run(); err != nil {
		http.Error(w, "Error running package.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Running package..."))
}

func ShutdownPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])


	checkIfPackageNameIsValid(pkgName, w)

	pkg, err := getPackage(pkgName)

	if err != nil {
		fmt.Println("Error:", err)
    return
	}

	if err := pkg.Shutdown(); err != nil {
		http.Error(w, "Error shutting down package process.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Shutting down package..."))
}

func InitPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	checkIfPackageNameIsValid(pkgName, w)

	pkg, err := getPackage(pkgName)

	if err != nil {
		fmt.Println("Error:", err)
    return
	}

	if err := pkg.Init(); err != nil {
		http.Error(w, "Error initializing package.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Initializing package..."))
}