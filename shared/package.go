package shared

import (
	"context"
	"os"
	"time"

	"google.golang.org/grpc"

	pc "rounds.com.ar/watcher/sdk/base/package-config"
	pb "rounds.com.ar/watcher/sdk/package"
)

type Package struct {
	Metadata  	*pc.PackageConfig

	Runtimepath string
	TempDir     string

	Client      pb.PackageServiceClient

	Process     *os.Process
	Connection  *grpc.ClientConn
	BasePort    int
}

func (h *Package) SetPorts(ports []int32) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	
	_, err := h.Client.SetPorts(ctx, &pb.SetPortsRequest{Ports: ports})
	return err == nil
}

func (h *Package) Install() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	
	resp, err := h.Client.Install(ctx, &pb.Empty{})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

func (h *Package) Run() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	
	resp, err := h.Client.Run(ctx, &pb.Empty{})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

func (h *Package) Exit() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
			
	_, err := h.Client.Exit(ctx, &pb.Empty{})
	return err
}