package shared

import (
	"context"
	"fmt"
	"time"

	pb "rounds.com.ar/sdk/shared/package"
)

// CallInstall calls the Install method on a package
func (h *PackageHost) CallInstall(packageName string) (bool, error) {
	for _, info := range h.Packages {
		if info.Name == packageName && info.Client != nil {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			
			resp, err := info.Client.Install(ctx, &pb.Empty{})
			if err != nil {
				return false, err
			}
			return resp.Success, nil
		}
	}
	return false, fmt.Errorf("package '%s' not found or not loaded", packageName)
}

// CallRun calls the Run method on a package
func (h *PackageHost) CallRun(packageName string) (bool, error) {
	for _, info := range h.Packages {
		if info.Name == packageName && info.Client != nil {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			
			resp, err := info.Client.Run(ctx, &pb.Empty{})
			if err != nil {
				return false, err
			}
			return resp.Success, nil
		}
	}
	return false, fmt.Errorf("package '%s' not found or not loaded", packageName)
}

// CallStop calls the Stop method on a package
func (h *PackageHost) CallExit(packageName string) error {
	for _, info := range h.Packages {
		if info.Name == packageName && info.Client != nil {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			
			_, err := info.Client.Exit(ctx, &pb.Empty{})
			return err
		}
	}
	return fmt.Errorf("package '%s' not found or not loaded", packageName)
}