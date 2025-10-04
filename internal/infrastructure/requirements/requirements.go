package requirements

import (
	"context"

	"github.com/rabbytesoftware/quiver/internal/models/requirement"
	"github.com/rabbytesoftware/quiver/internal/models/system"
)

type Requirements struct {
}

func NewRequirements() SRVInterface {
	return &Requirements{}
}

func (r *Requirements) Validate(
	ctx context.Context,
	requirements *requirement.Requirement,
) (bool, error) {
	return false, nil
}

func (r *Requirements) ValidateOS(
	ctx context.Context,
	recommendedOS system.OS,
) (bool, error) {
	return false, nil
}

func (r *Requirements) ValidateOSVersion(
	ctx context.Context,
	recommendedVersion string,
) (bool, error) {
	return false, nil
}

func (r *Requirements) ValidateArch(
	ctx context.Context,
	recommendedArch string,
) (bool, error) {
	return false, nil
}

func (r *Requirements) ValidateCPU(
	ctx context.Context,
	recommendedCPU int,
) (bool, error) {
	return false, nil
}

func (r *Requirements) ValidateMemory(
	ctx context.Context,
	recommendedMemory int,
) (bool, error) {
	return false, nil
}

func (r *Requirements) ValidateDisk(
	ctx context.Context,
	recommendedDisk int,
) (bool, error) {
	return false, nil
}

func (r *Requirements) ValidateNetwork(
	ctx context.Context,
	recommendedNetwork int,
) (bool, error) {
	return false, nil
}
