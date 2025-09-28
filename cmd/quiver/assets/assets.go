package assets

import (
	_ "embed"
)

//go:embed icons/app.ico
var WindowsIcon []byte

//go:embed icons/app.icns
var MacOSIcon []byte

//go:embed icons/app.png
var LinuxIcon []byte

//go:embed icons/app-256.png
var Icon256 []byte

//go:embed icons/app-128.png
var Icon128 []byte

//go:embed icons/app-64.png
var Icon64 []byte

//go:embed icons/app-32.png
var Icon32 []byte

//go:embed icons/app-16.png
var Icon16 []byte
