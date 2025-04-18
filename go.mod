module rounds.com.ar/watcher

go 1.22.5

require (
	github.com/gdamore/tcell/v2 v2.8.1
	github.com/pterm/pterm v0.12.80
	github.com/rivo/tview v0.0.0-20250330220935-949945f8d922
	rounds.com.ar/watcher/sdk v0.0.0
)

require (
	atomicgo.dev/cursor v0.2.0 // indirect
	atomicgo.dev/keyboard v0.2.9 // indirect
	atomicgo.dev/schedule v0.1.0 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/lithammer/fuzzysearch v1.1.8 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.org/x/term v0.28.0 // indirect
)

require (
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.4 // indirect
)

replace rounds.com.ar/watcher/sdk => ../shared
