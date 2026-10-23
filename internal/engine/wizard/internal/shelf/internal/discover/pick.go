package discover

import (
	"sort"
	"strings"
	"unicode"

	"github.com/rabbytesoftware/quiver.core/internal/domain"
	"github.com/rabbytesoftware/quiver.core/internal/engine/wizard/internal/shelf/internal/fsguard"
	"github.com/rabbytesoftware/quiver.core/internal/engine/wizard/internal/shelf/internal/models"
)

func Pick(
	found []models.Candidate,
	repo string,
	entryName string,
) (models.Candidate, string) {
	if len(found) == 0 {
		return models.Candidate{}, models.ReasonNoDesktop
	}
	if len(found) == 1 {
		return found[0], ""
	}

	matches := matching(found, repo, entryName)
	if len(matches) == 0 {
		return models.Candidate{}, models.ReasonAmbiguous
	}
	return matches[0], ""
}

func Renamed(
	c models.Candidate,
	name string,
) models.Candidate {
	if fsguard.SafeName(name) {
		c.Name = name
	}
	return c
}

func Icon(
	req models.Request,
	entry domain.ExposeEntry,
	c models.Candidate,
) string {
	icon := req.Media.Icon
	if c.Icon != "" {
		icon = c.Icon
	}
	if entry.Icon != "" {
		icon = fsguard.ExpandPath(entry.Icon, req.Workdir)
	}
	if strings.Contains(icon, "://") || fsguard.UnsafePath(icon, "") {
		return ""
	}
	return icon
}

func RepoName(
	bare domain.Namespace,
) string {
	parts := strings.Split(string(bare), domain.NamespaceSeparator)
	if len(parts) < 3 {
		return ""
	}
	return parts[2]
}

func matching(
	found []models.Candidate,
	repo string,
	entryName string,
) []models.Candidate {
	matches := []models.Candidate{}
	for _, c := range found {
		if sameName(c.Name, repo) || (entryName != "" && sameName(c.Name, entryName)) {
			matches = append(matches, c)
		}
	}
	sort.SliceStable(matches, func(i, j int) bool { return matches[i].Depth < matches[j].Depth })
	return matches
}

// sameName compares an app's name with a repository or entry name the way a
// person would: a repository is PrismLauncher where its bundle is Prism
// Launcher.app, and neither case nor separators tell two apps apart.
func sameName(
	a string,
	b string,
) bool {
	return strings.EqualFold(a, b) || (squash(a) != "" && squash(a) == squash(b))
}

func squash(
	name string,
) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case ' ', '-', '_', '.':
			return -1
		}
		return unicode.ToLower(r)
	}, name)
}

func shallowest(
	found []models.Candidate,
) []models.Candidate {
	depth := found[0].Depth
	for _, c := range found {
		depth = min(depth, c.Depth)
	}

	top := []models.Candidate{}
	for _, c := range found {
		if c.Depth == depth {
			top = append(top, c)
		}
	}
	return top
}

func uniqueNames(
	candidates []models.Candidate,
) []models.Candidate {
	seen := map[string]bool{}
	unique := []models.Candidate{}
	for _, c := range candidates {
		if seen[c.Name] {
			continue
		}
		seen[c.Name] = true
		unique = append(unique, c)
	}
	return unique
}
