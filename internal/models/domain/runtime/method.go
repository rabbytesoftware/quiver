package runtime

import "github.com/rabbytesoftware/quiver/internal/models/domain/system"

type Method struct {
	OS system.OS `json:"os"`
	Command []string `json:"command"`
}
