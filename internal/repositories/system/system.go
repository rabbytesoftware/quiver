package system

import (
	"github.com/rabbytesoftware/quiver/internal/core/metadata"
	"github.com/rabbytesoftware/quiver/internal/infrastructure"
)

type SystemRepository struct {
	infrastructure *infrastructure.Infrastructure
}

func NewSystemRepository(
	infrastructure *infrastructure.Infrastructure,
) SystemInterface {
	return &SystemRepository{
		infrastructure: infrastructure,
	}
}

func (s *SystemRepository) GetMetadata() *metadata.Metadata {
	return nil
}

func (s *SystemRepository) UpdateQuiver() error {
	return nil
}

func (s *SystemRepository) UninstallQuiver() error {
	return nil
}

func (s *SystemRepository) GetLogs() string {
	return ""
}

func (s *SystemRepository) RestartQuiver() error {
	return nil
}

func (s *SystemRepository) Status() string {
	return ""
}

func (s *SystemRepository) StopQuiver() error {
	return nil
}
