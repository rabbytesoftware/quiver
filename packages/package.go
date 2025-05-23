package packages

import (
	"context"
	"fmt"
	"time"

	pc "github.com/rabbytesoftware/quiver.sdk/base/package-config"
	pb "github.com/rabbytesoftware/quiver.sdk/package"
)

type Package struct {
	*pc.PackageConfig

	Runtime *PackageRuntime
}

/*
 * Init the child process
 */
func (pkg *Package) Init() error {
	err := pkg.extract()
	if err != nil {
		return err
	}

	err = pkg.Runtime.InitPackage()
	if err != nil {
		pkg.clean()
		return err
	}

	return nil
}

/*
 * Stop the child process
 */
func (pkg *Package) Shutdown() error {
	err := pkg.exit()
	if err != nil {
		return fmt.Errorf("failed to exit the package: %w", err)
	}

	err = pkg.Runtime.StopPackage()
	pkg.clean()

	return err
}

/*
 * Install the package
 */
func (pkg *Package) Install() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	
	resp, err := pkg.Runtime.Client.Install(ctx, &pb.Empty{})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

/*
 * Uninstall the package
 */
func (pkg *Package) Run() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	
	resp, err := pkg.Runtime.Client.Run(ctx, &pb.Empty{})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

/*
 * Set the ports for the package
 * @param ports []int32
 * @return bool
 */
func (pkg *Package) SetPorts(ports []int32) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	
	_, err := pkg.Runtime.Client.SetPorts(ctx, &pb.SetPortsRequest{Ports: ports})
	return err == nil
}
