package packages

import "strings"

type PackageNamespace string

func (p PackageNamespace) IsValid() bool {
	parts := strings.Split(string(p), "@")
	if len(parts) != 2 {
		return false
	}

	// Check that both parts are non-empty
	if parts[0] == "" || parts[1] == "" {
		return false
	}

	return true
}

func (p PackageNamespace) String() string {
	return string(p)
}
