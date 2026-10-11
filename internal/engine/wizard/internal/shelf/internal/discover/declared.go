package discover

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rabbytesoftware/quiver.core/internal/domain"
	"github.com/rabbytesoftware/quiver.core/internal/engine/wizard/internal/shelf/internal/fsguard"
	"github.com/rabbytesoftware/quiver.core/internal/engine/wizard/internal/shelf/internal/models"
	"github.com/rabbytesoftware/quiver.core/internal/engine/wizard/internal/workfs"
)

func Declared(
	req models.Request,
	entry domain.ExposeEntry,
) ([]models.Candidate, string, error) {
	if !fsguard.SafeName(entry.Name) {
		return nil, models.ReasonUnsafeName, nil
	}

	if name, rest, ok := domain.SplitExposeRoot(entry.Path); ok {
		return external(req, entry.Name, name, rest)
	}

	target := fsguard.ExpandPath(entry.Path, req.Workdir)
	if !workfs.Inside(req.Workdir, target) {
		return nil, models.ReasonOutsideWorkdir, nil
	}

	reason, err := fsguard.Contained(req.Workdir, target)
	if err != nil || reason != "" {
		return nil, reason, err
	}
	return []models.Candidate{{Name: entry.Name, Target: target, Declared: true}}, "", nil
}

// external resolves a desktop entry below an allowed external root. The result
// is only ever launched and linked to, never written to or removed.
func external(
	req models.Request,
	entryName string,
	root string,
	rest string,
) ([]models.Candidate, string, error) {
	dir, ok := req.Layout.External[root]
	if !ok || strings.Contains(rest, "..") || strings.Contains(rest, "${") {
		return nil, models.ReasonOutsideWorkdir, nil
	}
	target := filepath.Join(dir, filepath.FromSlash(strings.ReplaceAll(rest, `\`, "/")))
	if !fsguard.HasSuffixFold(target, ".exe") || fsguard.UnsafePath(target, `"%`) {
		return nil, models.ReasonUnsafePath, nil
	}
	if strings.ContainsAny(filepath.Base(target), "*?") {
		matches, err := filepath.Glob(target)
		if err != nil || len(matches) == 0 {
			return nil, models.ReasonNotFound, nil
		}
		if len(matches) > 1 {
			return nil, models.ReasonAmbiguous, nil
		}
		target = matches[0]
	}
	if _, ok := fsguard.ResolveInside(dir, target); !ok {
		if _, err := os.Lstat(target); err == nil {
			return nil, models.ReasonOutsideWorkdir, nil
		}
	}
	return []models.Candidate{{Name: entryName, Target: target, Declared: true}}, "", nil
}
