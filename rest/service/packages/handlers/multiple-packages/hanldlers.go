package multiple_packages_handlers

import (
	"encoding/json"
	"net/http"

	logger "github.com/rabbytesoftware/quiver/logger"
	packages "github.com/rabbytesoftware/quiver/packages"
)

type PackagesHandler struct {
	logs *logger.Logger
	pkgs *map[string]*packages.Package
}

func NewPackagesHandler(
	logs *logger.Logger,
	pkgs *map[string]*packages.Package,
) *PackagesHandler {
	return &PackagesHandler{
		logs: logs,
		pkgs: pkgs,
	}
}

func (h *PackagesHandler) GetPackagesList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader( http.StatusOK )
	json.NewEncoder(w).Encode( h.pkgs )
}