package shelf

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rabbytesoftware/quiver.core/internal/domain"
	"github.com/rabbytesoftware/quiver.core/internal/engine/wizard/internal/shelf/internal/models"
	"github.com/rabbytesoftware/quiver.core/internal/engine/wizard/internal/shelf/mocks"
)

func desktopExpose(
	path string,
) domain.Expose {
	return domain.Expose{Desktop: []domain.ExposeEntry{{Name: "Vendor", Path: path}}}
}

func TestShelf_Apply_ExternalDesktopApp(t *testing.T) {
	f := newFixture(t, "windows")
	programs := t.TempDir()
	f.Env["PROGRAMFILES"] = programs
	exe := filepath.Join(programs, "Vendor 3", "bin", "vendor-3.2.exe")
	mocks.WriteFile(t, exe, "x", 0o755)
	wd := f.Workdir(t, mocks.NsA)
	ctx := context.Background()

	testCases := []struct {
		name string
		path string
	}{
		{name: "exact file", path: `${PROGRAMFILES}/Vendor 3/bin/vendor-3.2.exe`},
		{name: "wildcard file name", path: `${PROGRAMFILES}\Vendor 3\bin\vendor-3*.exe`},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := f.shelf.Apply(ctx, mocks.NsA, wd, desktopExpose(tc.path), domain.ArrowMedia{})
			require.NoError(t, err)
			require.Empty(t, got.Refused)
			require.Len(t, got.Entries, 1)

			assert.Equal(t, exe, readShortcut(t, got.Entries[0].Location).Target, "the shortcut points at the external exe")
			app, err := f.shelf.AppEntry(ctx, wd)
			require.NoError(t, err)
			assert.Equal(t, exe, app)
		})
	}

	require.NoError(t, f.shelf.Remove(ctx, wd))
	assert.FileExists(t, exe, "uninstall must leave the vendor's app alone")
	assert.DirExists(t, filepath.Join(programs, "Vendor 3", "bin"))
	entries, err := os.ReadDir(lnkDir(filepath.Join(f.UserHome, "AppData", "Roaming")))
	require.NoError(t, err)
	assert.Empty(t, entries, "only quiver's own shortcut is removed")
}

func TestShelf_Apply_ExternalDesktopAppRefusals(t *testing.T) {
	f := newFixture(t, "windows")
	programs := t.TempDir()
	f.Env["PROGRAMFILES"] = programs
	outside := filepath.Join(t.TempDir(), "evil.exe")
	mocks.WriteFile(t, outside, "x", 0o755)
	mocks.WriteFile(t, filepath.Join(programs, "a", "one.exe"), "x", 0o755)
	mocks.WriteFile(t, filepath.Join(programs, "a", "two.exe"), "x", 0o755)
	mocks.WriteFile(t, filepath.Join(programs, "a", "tool.txt"), "x", 0o755)
	wd := f.Workdir(t, mocks.NsA)

	testCases := []struct {
		name   string
		path   string
		reason string
	}{
		{name: "traversal", path: `${PROGRAMFILES}/../evil.exe`, reason: models.ReasonOutsideWorkdir},
		{name: "unknown placeholder is not resolved", path: `${SYSTEMROOT}/evil.exe`, reason: models.ReasonNotFound},
		{name: "root not set on this host", path: `${LOCALAPPDATA}/a/one.exe`, reason: models.ReasonOutsideWorkdir},
		{name: "not an exe", path: `${PROGRAMFILES}/a/tool.txt`, reason: models.ReasonUnsafePath},
		{name: "ambiguous wildcard", path: `${PROGRAMFILES}/a/*.exe`, reason: models.ReasonAmbiguous},
		{name: "wildcard matching nothing", path: `${PROGRAMFILES}/a/nope*.exe`, reason: models.ReasonNotFound},
		{name: "missing file", path: `${PROGRAMFILES}/a/missing.exe`, reason: models.ReasonNotFound},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := f.shelf.Apply(context.Background(), mocks.NsA, wd, desktopExpose(tc.path), domain.ArrowMedia{})
			require.NoError(t, err)
			assert.Empty(t, got.Entries)
			require.Len(t, got.Refused, 1)
			assert.Equal(t, tc.reason, got.Refused[0].Reason)
		})
	}
}

func TestShelf_Apply_ExternalDesktopAppSymlinkEscape(t *testing.T) {
	mocks.RequireUnixHost(t)
	f := newFixture(t, "windows")
	programs := t.TempDir()
	f.Env["PROGRAMFILES"] = programs
	outside := filepath.Join(t.TempDir(), "evil.exe")
	mocks.WriteFile(t, outside, "x", 0o755)
	require.NoError(t, os.Symlink(outside, filepath.Join(programs, "link.exe")))
	wd := f.Workdir(t, mocks.NsA)

	got, err := f.shelf.Apply(context.Background(), mocks.NsA, wd, desktopExpose(`${PROGRAMFILES}/link.exe`), domain.ArrowMedia{})

	require.NoError(t, err)
	assert.Empty(t, got.Entries)
	require.Len(t, got.Refused, 1)
	assert.Equal(t, models.ReasonOutsideWorkdir, got.Refused[0].Reason)
}

func TestShelf_AppEntry_RejectsAMarkerOutsideEveryRoot(t *testing.T) {
	f := newFixture(t, "windows")
	f.Env["PROGRAMFILES"] = t.TempDir()
	outside := filepath.Join(t.TempDir(), "evil.exe")
	mocks.WriteFile(t, outside, "x", 0o755)
	wd := f.Workdir(t, mocks.NsA)
	mocks.WriteFile(t, filepath.Join(wd, appMarker), outside+"\n", 0o600)

	_, err := f.shelf.AppEntry(context.Background(), wd)

	assert.ErrorIs(t, err, ErrNoApp)
}
