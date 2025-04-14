package packages

import (
	"context"
	"fmt"
	"time"

	pc "rounds.com.ar/watcher/sdk/base/package-config"
	pb "rounds.com.ar/watcher/sdk/package"
)

type Package struct {
	*pc.PackageConfig

	Runtime *PackageRuntime
}

func (pkg *Package) Start() error {
	err := pkg.extract()
	if err != nil {
		return err
	}

	err = pkg.Runtime.StartPackage()
	if err != nil {
		pkg.clean()
		return err
	}

	return nil
}

func (pkg *Package) Remove() error {
	err := pkg.Exit()
	if err != nil {
		fmt.Println("Error exiting the package:", err)
	}

	err = pkg.Runtime.StopPackage()
	pkg.clean()

	return err
}

func (pkg *Package) SetPorts(ports []int32) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	
	_, err := pkg.Runtime.Client.SetPorts(ctx, &pb.SetPortsRequest{Ports: ports})
	return err == nil
}

func (pkg *Package) Install() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	
	resp, err := pkg.Runtime.Client.Install(ctx, &pb.Empty{})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

func (pkg *Package) Run() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	
	resp, err := pkg.Runtime.Client.Run(ctx, &pb.Empty{})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

func (pkg *Package) Exit() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
			
	_, err := pkg.Runtime.Client.Exit(ctx, &pb.Empty{})
	return err
}