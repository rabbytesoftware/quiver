package system

import "github.com/rabbytesoftware/quiver/internal/core/metadata"

type SystemInterface interface {
	GetMetadata() *metadata.Metadata
	UpdateQuiver() error
	UninstallQuiver() error

	GetLogs() string
	RestartQuiver() error
	Status() string
	StopQuiver() error
}