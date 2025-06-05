package packages

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/rabbytesoftware/quiver/logger"

	v0_1 "github.com/rabbytesoftware/quiver/packages/manifest/v0.1"
	yaml "gopkg.in/yaml.v3"
)

type ArrowsServer struct {
	Packages 	map[string]interface{}
	Repository 	string
	
	logs 		*logger.Logger
}

func NewArrowsServer(
	Repository string,
) *ArrowsServer {
	return &ArrowsServer{
		Repository: Repository,
		Packages:  	make(map[string]interface{}),
		logs:		logger.NewLogger("pkg"),
	}
}

func (as *ArrowsServer) Load(path string) error {
    // ? Read the YAML file
    data, err := ioutil.ReadFile(path)
    if err != nil {
        as.logs.Error("Failed to read YAML file: %v", err)
        return fmt.Errorf("failed to read YAML file: %w", err)
    }

    // ? First, parse just the version to determine which model to use
    var versionInfo struct {
        Version string `yaml:"version"`
    }

    if err := yaml.Unmarshal(
		data, 
		&versionInfo,
	); err != nil {
        as.logs.Error("Failed to parse version from YAML: %v", err)
        return fmt.Errorf("failed to parse version from YAML: %w", err)
    }

    if versionInfo.Version == "" {
        return fmt.Errorf("version field is required in YAML file")
    }

    // ? Load the package based on version
    var arrowObj interface{}

	// TODO: char2cs: This a horrible solution that copilot suggested.
	// TODO: 			We should have a better way to handle different versions of the Arrow struct...
	// TODO: 			Maybe use a factory pattern or a registry of version handlers? idk.
	// TODO: 			Really, how hard can a package manager be to build? 
	// TODO: 				Now I pray even more to the Nix gods...

    switch versionInfo.Version {
    case "0.1": // ? Create instance of v1.0 Arrow struct
        var arrow v0_1.Arrow
        if err := yaml.Unmarshal(data, &arrow); err != nil {
            as.logs.Error("Failed to unmarshal into v1.0 Arrow struct: %v", err)
            return fmt.Errorf("failed to unmarshal into v1.0 Arrow struct: %w", err)
        }
        arrowObj = &arrow
    default:
        return fmt.Errorf("unsupported version: %s", versionInfo.Version)
    }

    // ? Store the loaded package
    packageKey := fmt.Sprintf("%s @ %s ; with Arrow Manifest version: %s", as.Packages[path].Metadata().Name(), packageVersion, versionInfo.Version)
	as.Packages[path] = arrowObj

    as.logs.Info("Successfully loaded package: %s from %s", packageKey, filepath.Base(path))

    return nil
}
