package fetchnshare

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rabbytesoftware/quiver/internal/core/watcher"
)

type FNS struct {
}

func NewFNS() FNSInterface {
	fns := &FNS{}

	result, _ := fns.GetInfo(context.Background(), "C:/Users/Joaquin/Desktop/Code")

	watcher.Info(result.Path)
	watcher.Info(strconv.FormatInt(result.Size, 10))
	watcher.Info(result.ModTime.String())
	watcher.Info(string(result.Type))

	return fns
}

// GetInfo retrieves metadata information about a resource (file or directory).
// It returns ResourceInfo containing size, permissions, modification time, and other attributes.
// Supports both local filesystem paths and remote URLs (HTTP/HTTPS).
func (f *FNS) GetInfo(ctx context.Context, path string) (*ResourceInfo, error) {

	info := &ResourceInfo{Path: path}

	// Remote URLs
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {

		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)

		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch URL info: %w", err)
		}

		defer resp.Body.Close()

		info.Type = ResourceType(http.DetectContentType([]byte(path)))

		info.Size = resp.ContentLength // May be -1 if unknown
		if info.Size < 0 {

			var total int64 // var used in either Content-Range or full body read

			// If Content-Length doesn't exist, try Content-Range
			if cr := resp.Header.Get("Content-Range"); cr != "" {
				req.Header.Set("Range", "bytes=0-0")
				conRange := resp.Header.Get("Content-Range")

				_, err := fmt.Sscanf(conRange, "bytes 0-0/%d", &total)
				if err == nil {
					info.Size = total
				}
			} else {
				// As a last resort, read the entire body to determine size
				var total int64
				buf := make([]byte, 32*1024)
				for {
					n, err := resp.Body.Read(buf)
					total += int64(n)
					if err == io.EOF {
						break
					}
				}
				info.Size = total
			}
		}

		tim := resp.Header.Get("Last-Modified")
		modT, err := http.ParseTime(tim)
		if tim == "" || err != nil { // If parsing fails or header is missing
			info.ModTime = time.Time{} // Unknown mod time (or could use time.Now())
		} else {
			info.ModTime = modT
		}
	} else { // Local filesystem paths

		stat, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("failed to stat path: %w", err)
		}

		info.Size = stat.Size()
		info.ModTime = stat.ModTime()
		if stat.IsDir() { // Maybe check for more types?
			info.Type = ResourceType("directory")
		} else {
			info.Type = ResourceType("file")
		}
	}

	return info, nil
}

// Exists checks whether a resource exists at the given path.
// Returns true if the resource exists, false otherwise.
// Works with both local filesystem paths and remote URLs.
func (f *FNS) Exists(ctx context.Context, path string) (bool, error) {
	return false, nil
}

// IsDir checks whether the resource at the given path is a directory.
// Returns true if the resource is a directory, false if it's a file or doesn't exist.
// Only works with local filesystem paths.
func (f *FNS) IsDir(ctx context.Context, path string) (bool, error) {
	return false, nil
}

// IsFile checks whether the resource at the given path is a regular file.
// Returns true if the resource is a file, false if it's a directory or doesn't exist.
// Only works with local filesystem paths.
func (f *FNS) IsFile(ctx context.Context, path string) (bool, error) {
	return false, nil
}

// Read reads the entire content of a resource into memory as a byte slice.
// Returns the complete file content or downloaded data.
// Use ReadStream for large files to avoid memory issues.
func (f *FNS) Read(ctx context.Context, path string) ([]byte, error) {
	return nil, nil
}

// ReadStream returns an io.ReadCloser for streaming data from a resource.
// Preferred for large files as it doesn't load everything into memory.
// Caller must close the returned ReadCloser when done.
func (f *FNS) ReadStream(ctx context.Context, path string) (io.ReadCloser, error) {
	return nil, nil
}

// Write writes data to a resource, creating the file if it doesn't exist.
// Overwrites existing files. Only works with local filesystem paths.
// Use WriteStream for large data to avoid memory issues.
func (f *FNS) Write(ctx context.Context, path string, data []byte) error {
	return nil
}

// WriteStream writes data from an io.Reader to a resource.
// Preferred for large data as it streams without loading everything into memory.
// Only works with local filesystem paths.
func (f *FNS) WriteStream(ctx context.Context, path string, reader io.Reader) error {
	return nil
}

// Append appends data to the end of a resource, creating the file if it doesn't exist.
// Only works with local filesystem paths.
func (f *FNS) Append(ctx context.Context, path string, data []byte) error {
	return nil
}

// List returns a slice of ResourceInfo for all items in a directory.
// Only works with local filesystem paths.
func (f *FNS) List(ctx context.Context, path string) ([]ResourceInfo, error) {
	return nil, nil
}

// Mkdir creates a single directory with the specified permissions.
// Fails if parent directories don't exist. Only works with local filesystem paths.
func (f *FNS) Mkdir(ctx context.Context, path string, perm os.FileMode) error {
	return nil
}

// MkdirAll creates a directory and all necessary parent directories with the specified permissions.
// Creates parent directories as needed. Only works with local filesystem paths.
func (f *FNS) MkdirAll(ctx context.Context, path string, perm os.FileMode) error {
	return nil
}

// Remove deletes a single file or empty directory.
// Fails if the directory is not empty. Only works with local filesystem paths.
func (f *FNS) Remove(ctx context.Context, path string) error {
	return nil
}

// RemoveAll deletes a file or directory and all its contents recursively.
// Use with caution as it permanently deletes everything. Only works with local filesystem paths.
func (f *FNS) RemoveAll(ctx context.Context, path string) error {
	return nil
}

// Copy copies a resource from source to destination.
// Works with local files and can download from URLs to local destinations.
func (f *FNS) Copy(ctx context.Context, src, dst string) error {
	return nil
}

// Move moves a resource from source to destination.
// Equivalent to copy + remove. Only works with local filesystem paths.
func (f *FNS) Move(ctx context.Context, src, dst string) error {
	return nil
}

// Rename renames a resource from source to destination.
// Alias for Move. Only works with local filesystem paths.
func (f *FNS) Rename(ctx context.Context, src, dst string) error {
	return nil
}

// Chmod changes the file permissions of a resource.
// Only works with local filesystem paths.
func (f *FNS) Chmod(ctx context.Context, path string, mode os.FileMode) error {
	return nil
}

// Chown changes the ownership of a resource (user ID and group ID).
// Only works with local filesystem paths and requires appropriate permissions.
func (f *FNS) Chown(ctx context.Context, path string, uid, gid int) error {
	return nil
}

// Download downloads a resource from a URL to a local destination path.
// The progress callback receives the number of bytes downloaded.
func (f *FNS) Download(ctx context.Context, url, dst string, progress func(int)) error {
	return nil
}

// DownloadStream returns an io.ReadCloser for streaming a download from a URL.
// The progress callback receives the number of bytes downloaded.
// Caller must close the returned ReadCloser when done.
func (f *FNS) DownloadStream(ctx context.Context, url string, progress func(int)) (io.ReadCloser, error) {
	return nil, nil
}

// Fetch downloads content from a URL and returns it as a byte slice.
// Use for small resources. For large downloads, use DownloadStream.
func (f *FNS) Fetch(ctx context.Context, url string) ([]byte, error) {
	return nil, nil
}

// CacheGet retrieves data from the cache using the specified key.
// Returns an error if the key doesn't exist or has expired.
func (f *FNS) CacheGet(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

// CacheSet stores data in the cache with the specified key and time-to-live (TTL).
// Data will be automatically removed after the TTL expires.
func (f *FNS) CacheSet(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	return nil
}

// CacheDelete removes data from the cache using the specified key.
// No error is returned if the key doesn't exist.
func (f *FNS) CacheDelete(ctx context.Context, key string) error {
	return nil
}

// CacheClear removes all data from the cache.
// Use with caution as this permanently deletes all cached data.
func (f *FNS) CacheClear(ctx context.Context) error {
	return nil
}

// Resolve determines the actual path and resource type for a given path or URL.
// Returns the resolved path, resource type, and any error encountered.
func (f *FNS) Resolve(ctx context.Context, path string) (string, ResourceType, error) {
	return "", "", nil
}

// Validate checks if a path or URL is valid and accessible.
// Returns an error if the resource is invalid, inaccessible, or blocked by policy.
func (f *FNS) Validate(ctx context.Context, path string) error {
	return nil
}

// TempFile creates a temporary file with the specified pattern and returns its path.
// The file is created in the system's temporary directory.
func (f *FNS) TempFile(ctx context.Context, pattern string) (string, error) {
	return "", nil
}

// TempDir creates a temporary directory with the specified pattern and returns its path.
// The directory is created in the system's temporary directory.
func (f *FNS) TempDir(ctx context.Context, pattern string) (string, error) {
	return "", nil
}

// Walk recursively traverses a directory tree, calling the provided function for each file and directory.
// The callback function receives the path, ResourceInfo, and any error encountered.
// Return an error from the callback to stop the walk, or nil to continue.
func (f *FNS) Walk(ctx context.Context, root string, fn func(path string, info ResourceInfo, err error) error) error {
	return nil
}
