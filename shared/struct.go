package shared

import (
	"os"
	"sync"

	"google.golang.org/grpc"

	pb "rounds.com.ar/sdk/shared/package"
)

// PackageInfo contains basic information about a package
type PackageInfo struct {
	PackageName string
	Version     string

	Icon 		string
	Name        string
	Description string

	Process     *os.Process
	Client      pb.PackageServiceClient
	Connection  *grpc.ClientConn
	BasePort    int
}

// PackageHost manages the lifecycle of packages
type PackageHost struct {
	PackagesDir string
	Packages    map[string]*PackageInfo
	NextPort    int
	mutex       sync.Mutex
}
