package manifest

import (
	v0_1 "github.com/rabbytesoftware/quiver/packages/manifest/v0.1"
)

type ArrowInterface interface {
	// Basic information
	Manifest() string
	Name() string
	Description() string
	Mainteiners() []string
	Credits() string
	License() string
	Repository() string
	Documentation() string
	ArrowVersion() string

	// Access methods for internal structures
	GetMetadata() *v0_1.Metadata
	GetRequirements() v0_1.Requirements
	GetDependencies() []string
	GetNetbridge() []v0_1.Port
	GetVariables() []v0_1.Variable
	GetMethods() v0_1.Methods
}

type Metadata interface {
	GetName() string
	GetDescription() string
	GetMaintainers() []string
	GetCredits() []string
	GetLicense() string
	GetRepository() string
	GetDocumentation() string
	GetVersion() string
}

type Requirements interface {
	GetMinimum() Requirement
	GetRecommended() Requirement
}

type Requirement interface {
	GetCpuCores() int
	GetRamGB() int
	GetDiskGB() int
	GetNetworkMbps() int
}

type Variable interface {
	GetName() string
	GetDefault() interface{}
	GetValues() []string
	GetMin() *int
	GetMax() *int
	GetSensitive() bool
}

type Netbridge interface {
	GetName() string
	GetProtocol() string
}

type Methods interface {
	GetInstall() map[string][]string
	GetExecute() map[string][]string
	GetUninstall() map[string][]string
	GetUpdate() map[string][]string
	GetValidate() map[string][]string
	GetMethod(methodName string) map[string][]string
}
