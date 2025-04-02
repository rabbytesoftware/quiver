package shared

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "rounds.com.ar/sdk/shared/package"
)

// startPackage launches a package process and connects to it
func (h *PackageHost) startPackage(path string, info *PackageInfo) error {
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	

	// =========================
	// 			Getters
	// =========================
	
	packageName, err := client.GetPackageName(ctx, &pb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to get package name: %w", err)
	}
	info.PackageName = packageName.Value

	versionResp, err := client.GetVersion(ctx, &pb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to get package version: %w", err)
	}
	info.Version = versionResp.Value

	nameResp, err := client.GetName(ctx, &pb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to get package name: %w", err)
	}
	info.Name = nameResp.Value

	iconResp, err := client.GetIcon(ctx, &pb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to get package version: %w", err)
	}
	info.Icon = iconResp.Value

	descResp, err := client.GetDescription(ctx, &pb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to get package version: %w", err)
	}
	info.Description = descResp.Value
	
	return nil
}