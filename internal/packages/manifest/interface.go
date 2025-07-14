package manifest

// ArrowInterface defines the main interface for arrow manifests
// This interface is version-agnostic and doesn't depend on specific version types
type ArrowInterface interface {
	// Basic information
	Manifest() string
	Name() string
	Description() string
	Maintainers() []string
	Credits() string
	License() string
	Repository() string
	Documentation() string
	ArrowVersion() string
	GetSupportedArchs(os string) []string
	
	// Access methods for internal structures (version-agnostic)
	GetMetadata() Metadata
	GetRequirements() Requirements
	GetDependencies() []string
	GetNetbridge() []Netbridge
	GetVariables() []Variable
	GetMethods() Methods
}

// Metadata interface for arrow metadata
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

// Requirements interface for arrow requirements
type Requirements interface {
	GetMinimum() Requirement
	GetRecommended() Requirement
	GetCompatible() map[string][]string
}

// Requirement interface for individual requirement specifications
type Requirement interface {
	GetCpuCores() int
	GetRamGB() int
	GetDiskGB() int
	GetNetworkMbps() int
}

// Variable interface for arrow variables
type Variable interface {
	GetName() string
	GetDefault() interface{}
	GetValues() []string
	GetMin() *int
	GetMax() *int
	GetSensitive() bool
}

// Netbridge interface for network bridge configurations
type Netbridge interface {
	GetName() string
	GetProtocol() string
}

// Methods interface for arrow execution methods
type Methods interface {
	GetInstall() map[string]map[string][]string
	GetExecute() map[string]map[string][]string
	GetUninstall() map[string]map[string][]string
	GetUpdate() map[string]map[string][]string
	GetValidate() map[string]map[string][]string
	GetMethod(methodName string) map[string]map[string][]string
}

// VersionParser interface for parsing version information
type VersionParser interface {
	ParseVersion(data []byte) (string, error)
}

// ArrowFactory interface for creating arrows from different versions
type ArrowFactory interface {
	CreateArrow(version string, data []byte) (ArrowInterface, error)
	SupportedVersions() []string
}
