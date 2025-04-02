package shared

import (
	"context"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "rounds.com.ar/sdk/shared/package"
)

type Package struct {
	Version     string
	Icon 		string
	Name        string
	Description string
	MaxPorts	int32

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