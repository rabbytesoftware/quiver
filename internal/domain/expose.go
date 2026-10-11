package domain

import "strings"

const ExposeAuto = "auto"

type ExposeKind string

const (
	ExposeKindCLI     ExposeKind = "cli"
	ExposeKindDesktop ExposeKind = "desktop"
)

type Expose struct {
	CLI     []ExposeEntry `yaml:"cli"     json:"cli,omitempty"`
	Desktop []ExposeEntry `yaml:"desktop" json:"desktop,omitempty"`
}

func (e Expose) IsEmpty() bool {
	return len(e.CLI) == 0 && len(e.Desktop) == 0
}

type ExposeEntry struct {
	Name       string   `yaml:"name"       json:"name"`
	Path       string   `yaml:"path"       json:"path"`
	Icon       string   `yaml:"icon"       json:"icon,omitempty"`
	Categories []string `yaml:"categories" json:"categories,omitempty"`
}

// ExposeExternalRoots are the only places outside its workdir a Windows desktop
// entry may point at, written as ${NAME} at the start of the path. They are
// the places vendor installers put an app, which Quiver never writes to.
var ExposeExternalRoots = []string{"LOCALAPPDATA", "PROGRAMFILES", "PROGRAMFILES(X86)", "APPDATA"}

// SplitExposeRoot splits a path that starts with an allowed external root into
// the root's name and the path below it.
func SplitExposeRoot(path string) (string, string, bool) {
	for _, root := range ExposeExternalRoots {
		rest, ok := strings.CutPrefix(path, "${"+root+"}")
		if ok && len(rest) > 1 && (rest[0] == '/' || rest[0] == '\\') {
			return root, rest[1:], true
		}
	}
	return "", "", false
}
