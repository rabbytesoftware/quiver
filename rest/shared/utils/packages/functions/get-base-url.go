package packages_functions

import (
	"os"
	"path/filepath"
)

func GetBaseUrl() string{
	exePath, err := os.Executable()

	if(err != nil){
		return ""
	}
	
	return filepath.Dir(exePath)
}

func GetPath(args ...string) string{
	base := GetBaseUrl()

	if(base == ""){
		return ""
	}

	all := append([]string{base}, args...)

	return filepath.Join(all...)
}