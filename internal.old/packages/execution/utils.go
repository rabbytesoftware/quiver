package execution

import (
	"os"
	"path/filepath"
	"strings"
)

// isPathSafe checks if a path is safe for extraction (prevents directory traversal)
func isPathSafe(path, workDir string) bool {
	cleanPath := filepath.Clean(path)
	cleanWorkDir := filepath.Clean(workDir)
	return strings.HasPrefix(cleanPath, cleanWorkDir+string(os.PathSeparator)) || cleanPath == cleanWorkDir
} 