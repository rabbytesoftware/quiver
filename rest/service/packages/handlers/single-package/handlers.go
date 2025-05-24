package single_package_handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	logger "github.com/rabbytesoftware/quiver/logger"
	packages "github.com/rabbytesoftware/quiver/packages"
)

type PackageHandler struct {
	logs *logger.Logger
	pkgs *map[string]*packages.Package
}

func NewPackageHandler(
	logs *logger.Logger,
	pkgs *map[string]*packages.Package,
) *PackageHandler {
	return &PackageHandler{
		logs: logs,
		pkgs: pkgs,
	}
}

// * Utils *
func (h *PackageHandler) checkIfPackageNameIsValid(name string, w http.ResponseWriter) {
	if len(name) == 0 {
		http.Error(w, "Package name is required.", http.StatusBadRequest)
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

// * Handlers *
func (h *PackageHandler) GetSinglePackageHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	h.checkIfPackageNameIsValid(pkgName, w)

	pkg := (*h.pkgs)[pkgName]
	if pkg == nil {
		http.Error(w, "Package not found.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pkg)
}

func (h *PackageHandler) DeletePackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	h.checkIfPackageNameIsValid(pkgName, w)

	// TODO
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

func (h *PackageHandler) InstallPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	h.checkIfPackageNameIsValid(pkgName, w)

	pkg := (*h.pkgs)[pkgName]
	if pkg == nil {
		http.Error(w, "Package not found.", http.StatusNotFound)
		return
	}

	if _, err := pkg.Install(); err != nil {
		http.Error(w, "An error occurred while installing the package.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("The package has been successfully installed."))
}

func (h *PackageHandler) RunPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	h.checkIfPackageNameIsValid(pkgName, w)

	pkg := (*h.pkgs)[pkgName]
	if pkg == nil {
		http.Error(w, "Package not found.", http.StatusNotFound)
		return
	}

	if _, err := pkg.Run(); err != nil {
		http.Error(w, "Error running package.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Running package..."))
}

func (h *PackageHandler) ShutdownPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])


	h.checkIfPackageNameIsValid(pkgName, w)

	pkg := (*h.pkgs)[pkgName]
	if pkg == nil {
		http.Error(w, "Package not found.", http.StatusNotFound)
		return
	}

	if err := pkg.Shutdown(); err != nil {
		http.Error(w, "Error shutting down package process.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Shutting down package..."))
}

func (h *PackageHandler) InitPackageHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	pkgName := strings.ToLower(params["name"])

	h.checkIfPackageNameIsValid(pkgName, w)

	pkg := (*h.pkgs)[pkgName]
	if pkg == nil {
		http.Error(w, "Package not found.", http.StatusNotFound)
		return
	}

	if err := pkg.Init(); err != nil {
		http.Error(w, "Error initializing package.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Initializing package..."))
}
