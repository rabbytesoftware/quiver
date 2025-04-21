package multiple_packages_handlers

import (
	"encoding/json"
	"net/http"

	packages_global_variables "rounds.com.ar/watcher/rest/shared/utils/packages/global-variables"
)

func GetPackagesList(w http.ResponseWriter, r *http.Request) {
	packages := packages_global_variables.Packages

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(packages)
}