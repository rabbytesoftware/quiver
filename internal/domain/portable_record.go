package domain

const PortableRecordFile = ".quiver-apps.json"

// ResourceForkDir is the directory macOS archivers add to a zip to carry
// resource forks and Finder metadata. It mirrors the real tree with AppleDouble
// stubs, so a bundle inside it is not an app and is never recorded or exposed.
const ResourceForkDir = "__MACOSX"

type PortableApp struct {
	Name  string `json:"name"           yaml:"name"`
	Entry string `json:"entry"          yaml:"entry"`
	Icon  string `json:"icon,omitempty" yaml:"icon,omitempty"`
}

type PortableRecord struct {
	Apps []PortableApp `json:"apps" yaml:"apps"`
}

func (r PortableRecord) Merge(
	apps []PortableApp,
) PortableRecord {
	incoming := make(map[string]PortableApp, len(apps))
	names := make(map[string]bool, len(apps))
	for _, app := range apps {
		incoming[app.Entry] = app
		names[app.Name] = true
	}

	consumed := make(map[string]bool, len(apps))
	merged := make([]PortableApp, 0, len(r.Apps)+len(apps))
	for _, existing := range r.Apps {
		replacement, ok := incoming[existing.Entry]
		if ok {
			merged = append(merged, replacement)
			consumed[existing.Entry] = true
			continue
		}
		if !names[existing.Name] {
			merged = append(merged, existing)
		}
	}

	for _, app := range apps {
		if consumed[app.Entry] {
			continue
		}
		merged = append(merged, app)
		consumed[app.Entry] = true
	}

	return PortableRecord{Apps: merged}
}
