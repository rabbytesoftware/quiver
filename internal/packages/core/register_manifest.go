package core

import (
	v0_1 "github.com/rabbytesoftware/quiver/internal/packages/manifest/v0.1"
	v0_2 "github.com/rabbytesoftware/quiver/internal/packages/manifest/v0.2"
)

func registerManifests() {
	v0_1.Init()
	v0_2.Init()
}
