package shared

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "rounds.com.ar/watcher/sdk/package"
	pkg "rounds.com.ar/watcher/shared"
)

// startPackage launches a package process and connects to it
func (h *PackagesHost) startPackage(path string, info *pkg.Package) error {
	// Launch the package with the port as an argument
	portStr := strconv.Itoa(info.BasePort)
	cmd := exec.Command(path, portStr)
	
	// Set up pipes for stderr to capture package output
	cmd.Stderr = os.Stderr
	
	// Start the package
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start package process: %w", err)
	}
	
	info.Process = cmd.Process
	
	// Wait for the gRPC server to start
	time.Sleep(1 * time.Second)
	
	// Connect to the package
	addr := fmt.Sprintf("localhost:%d", info.BasePort)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return fmt.Errorf("failed to connect to package: %w", err)
	}
	
	// Create a client
	client := pb.NewPackageServiceClient(conn)
	info.Client = client
	info.Connection = conn
	
	// Get package information

	return nil
}