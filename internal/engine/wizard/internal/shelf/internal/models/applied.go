package models

import (
	"github.com/rabbytesoftware/quiver.core/internal/domain"
)

type AppliedEntry struct {
	Kind     domain.ExposeKind
	Index    int
	Name     string
	Target   string
	Location string
}

type Applied struct {
	Entries []AppliedEntry
	Refused []Refusal
	Skipped []Refusal
	// Unclaimed are desktop entries refused only because the app already sits
	// unmanaged where it would be placed; the arrow's own copy still runs.
	Unclaimed []AppliedEntry
}

func (a *Applied) Record(
	kind domain.ExposeKind,
	index int,
	c Candidate,
	p Placement,
) {
	if p.Refused != "" {
		a.Refuse(kind, index, c.Name, p.Refused)
		if kind == domain.ExposeKindDesktop && p.Refused == ReasonUnmanaged {
			a.Unclaimed = append(a.Unclaimed, AppliedEntry{Kind: kind, Index: index, Name: c.Name, Target: c.Target})
		}
		return
	}
	a.Entries = append(a.Entries, AppliedEntry{
		Kind:     kind,
		Index:    index,
		Name:     c.Name,
		Target:   c.Target,
		Location: p.Location,
	})
}

func (a *Applied) Locations() map[string]bool {
	locations := map[string]bool{}
	for _, e := range a.Entries {
		locations[e.Location] = true
	}
	return locations
}

func (a *Applied) Skip(
	kind domain.ExposeKind,
	index int,
	name string,
	reason string,
) {
	a.Skipped = append(a.Skipped, Refusal{Kind: kind, Index: index, Name: name, Reason: reason})
}

func (a *Applied) Refuse(
	kind domain.ExposeKind,
	index int,
	name string,
	reason string,
) {
	a.Refused = append(a.Refused, Refusal{Kind: kind, Index: index, Name: name, Reason: reason})
}
