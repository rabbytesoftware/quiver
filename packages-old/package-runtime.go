package packages

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/rabbytesoftware/quiver.compiler/shared/package"
)

type PackageRuntime struct {
	PackagePath string
	RuntimePath string
	TempDir     string

	Client      pb.PackageServiceClient

	Process     *os.Process
	Connection  *grpc.ClientConn
	BasePort    int
}

// startPackage launches a package process and connects to it
func (pkg *PackageRuntime) InitPackage() error {
	// Launch the package with the port as an argument
	cmd := exec.Command(
		pkg.RuntimePath, 
		strconv.Itoa(pkg.BasePort),
	)
	
	// Set up pipes for stderr to capture package output
	cmd.Stderr = os.Stderr
	
	// Start the package
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start package process: %w", err)
	}
	
	pkg.Process = cmd.Process
	
	// Wait for the gRPC server to start
	time.Sleep(1 * time.Second)
	
	// Connect to the package
	addr := fmt.Sprintf("localhost:%d", pkg.BasePort)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(5*time.Second)) // TODO...
	if err != nil {
		return fmt.Errorf("failed to connect to package: %w", err)
	}
	
	// Create a client
	client := pb.NewPackageServiceClient(conn)
	pkg.Client = client
	pkg.Connection = conn

	return nil
}

func (pkg *PackageRuntime) StopPackage() error {
	if pkg.Process != nil {
		if err := pkg.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill package process: %w", err)
		}
		pkg.Process = nil
	}

	if pkg.Connection != nil {
		if err := pkg.Connection.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
		pkg.Connection = nil
	}

	return nil
}