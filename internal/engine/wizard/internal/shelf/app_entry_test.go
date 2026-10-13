package shelf

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rabbytesoftware/quiver.core/internal/domain"
	"github.com/rabbytesoftware/quiver.core/internal/engine/wizard/internal/shelf/internal/models"
	"github.com/rabbytesoftware/quiver.core/internal/engine/wizard/internal/shelf/mocks"
)

func TestShelf_AppEntry_FollowsTheDesktopEntryApplyPlaced(t *testing.T) {
	mocks.RequireUnixHost(t)
	f := newFixture(t, "linux")
	wd := f.Workdir(t, mocks.NsA)
	mocks.WriteFile(t, filepath.Join(wd, "Tool.AppImage"), "x", 0o755)
	mocks.WriteFile(t, filepath.Join(wd, "cli"), "x", 0o755)
	ctx := context.Background()

	_, err := f.shelf.Apply(ctx, mocks.NsA, wd, cliExpose("cli", "cli"), domain.ArrowMedia{})
	require.NoError(t, err)
	_, err = f.shelf.AppEntry(ctx, wd)
	assert.ErrorIs(t, err, ErrNoApp, "a CLI-only arrow has no app to start")

	desktop := domain.Expose{Desktop: []domain.ExposeEntry{{Name: "Tool", Path: "Tool.AppImage"}}}
	_, err = f.shelf.Apply(ctx, mocks.NsA, wd, desktop, domain.ArrowMedia{})
	require.NoError(t, err)
	got, err := f.shelf.AppEntry(ctx, wd)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(wd, "Tool.AppImage"), got)

	_, err = f.shelf.Apply(ctx, mocks.NsA, wd, cliExpose("cli", "cli"), domain.ArrowMedia{})
	require.NoError(t, err)
	_, err = f.shelf.AppEntry(ctx, wd)
	assert.ErrorIs(t, err, ErrNoApp, "an update that drops the desktop entry drops the app")
}

func TestShelf_AppEntry_RejectsWhatCannotBeStarted(t *testing.T) {
	mocks.RequireUnixHost(t)
	f := newFixture(t, "linux")
	wd := f.Workdir(t, mocks.NsA)
	ctx := context.Background()
	record := func(target string) {
		mocks.WriteFile(t, filepath.Join(wd, appMarker), target+"\n", 0o600)
	}

	record(wd)
	_, err := f.shelf.AppEntry(ctx, wd)
	assert.ErrorIs(t, err, ErrNoApp, "a directory")

	plain := filepath.Join(wd, "plain")
	mocks.WriteFile(t, plain, "x", 0o644)
	record(plain)
	_, err = f.shelf.AppEntry(ctx, wd)
	assert.ErrorIs(t, err, ErrNoApp, "no exec bit")

	record(filepath.Join(wd, "..", "..", "escape"))
	_, err = f.shelf.AppEntry(ctx, wd)
	assert.ErrorIs(t, err, ErrNoApp, "outside the arrow")

	record(filepath.Join(wd, "gone"))
	_, err = f.shelf.AppEntry(ctx, wd)
	assert.ErrorIs(t, err, ErrNoApp, "missing")
}

func TestBundleExecutable(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "Tool.app", "Contents", "MacOS")
	mocks.WriteFile(t, filepath.Join(dir, "Helper"), "x", 0o755)
	bundle := filepath.Dir(filepath.Dir(dir))

	_, err := bundleExecutable(bundle)
	assert.NoError(t, err, "a single executable is the app")

	mocks.WriteFile(t, filepath.Join(dir, "Other"), "x", 0o755)
	_, err = bundleExecutable(bundle)
	assert.Error(t, err, "two executables with none named like the bundle are ambiguous")

	mocks.WriteFile(t, filepath.Join(dir, "Tool"), "x", 0o755)
	got, err := bundleExecutable(bundle)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "Tool"), got, "the executable named like the bundle wins")
}

func TestShelf_AppEntry_StartsTheArrowsCopyWhenTheAppExistsUnmanaged(t *testing.T) {
	f := newFixture(t, "darwin")
	wd := f.Workdir(t, mocks.NsA)
	mocks.WriteBundle(t, filepath.Join(wd, "Tool.app"), "v1")
	mocks.WriteFile(t, filepath.Join(wd, "Tool.app", "Contents", "MacOS", "Tool"), "x", 0o755)
	existing := filepath.Join(f.Apps[0], "Tool.app")
	mocks.WriteBundle(t, existing, "user install")
	expose := domain.Expose{Desktop: []domain.ExposeEntry{{Name: "Tool", Path: domain.ExposeAuto}}}
	ctx := context.Background()

	got, err := f.shelf.Apply(ctx, mocks.NsA, wd, expose, domain.ArrowMedia{})
	require.NoError(t, err)
	require.Len(t, got.Refused, 1)
	assert.Equal(t, models.ReasonUnmanaged, got.Refused[0].Reason)

	app, err := f.shelf.AppEntry(ctx, wd)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(wd, "Tool.app", "Contents", "MacOS", "Tool"), app)

	require.NoError(t, f.shelf.Remove(ctx, wd))
	assert.Equal(t, "user install", mocks.ReadBundleVersion(t, existing), "an app Quiver never owned is left intact")
}
