package execution

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// handleSpecialCommand processes special commands like GET, UNCOMPRESS, MOVE, REMOVE
func (e *Engine) handleSpecialCommand(ctx context.Context, command, workDir string, env []string) (bool, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false, nil
	}
	
	cmd := strings.ToUpper(parts[0])
	
	e.logger.Debug("Checking for special command: %s", cmd)
	
	switch {
	case strings.HasPrefix(cmd, "GET:"):
		e.logger.Info("Executing special command: GET (HTTP Download)")
		return true, e.handleGetCommand(ctx, command, workDir)
	case strings.HasPrefix(cmd, "UNCOMPRESS:"):
		e.logger.Info("Executing special command: UNCOMPRESS (Archive Extraction)")
		return true, e.handleUncompressCommand(command, workDir)
	case strings.HasPrefix(cmd, "MOVE:"):
		e.logger.Info("Executing special command: MOVE (File Operation)")
		return true, e.handleMoveCommand(command, workDir)
	case strings.HasPrefix(cmd, "REMOVE:"):
		e.logger.Info("Executing special command: REMOVE (File Deletion)")
		return true, e.handleRemoveCommand(command, workDir)
	default:
		e.logger.Debug("Not a special command, falling back to shell execution")
		return false, nil
	}
}

// handleGetCommand implements HTTP downloads with verbose logging and progress tracking
func (e *Engine) handleGetCommand(ctx context.Context, command, workDir string) error {
	e.logger.Info("Starting HTTP download operation")
	
	// Parse GET: <URL>
	parts := strings.SplitN(command, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid GET command format: %s", command)
	}
	
	url := strings.TrimSpace(parts[1])
	if url == "" {
		return fmt.Errorf("empty URL in GET command")
	}
	
	e.logger.Info("Download URL: %s", url)
	e.logger.Info("Working directory: %s", workDir)
	
	// Create request with context for timeout/cancellation
	e.logger.Debug("Creating HTTP request...")
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		e.logger.Error("Failed to create HTTP request: %v", err)
		return fmt.Errorf("failed to create request: %v", err)
	}
	
	e.logger.Info("Sending HTTP request...")
	resp, err := e.httpClient.Do(req)
	if err != nil {
		e.logger.Error("HTTP request failed: %v", err)
		return fmt.Errorf("failed to download %s: %v", url, err)
	}
	defer resp.Body.Close()
	
	e.logger.Info("HTTP Response: %d %s", resp.StatusCode, resp.Status)
	
	if resp.StatusCode != http.StatusOK {
		e.logger.Error("HTTP error: %d %s", resp.StatusCode, resp.Status)
		return fmt.Errorf("HTTP error %d downloading %s", resp.StatusCode, url)
	}
	
	// Log response headers for debugging
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		e.logger.Info("Content-Length: %s bytes", contentLength)
	}
	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		e.logger.Info("Content-Type: %s", contentType)
	}
	
	// Extract filename from URL
	filename := filepath.Base(url)
	if filename == "." || filename == "/" {
		filename = "download"
	}
	
	// Save to working directory
	outputPath := filepath.Join(workDir, filename)
	e.logger.Info("Saving to: %s", outputPath)
	
	outFile, err := os.Create(outputPath)
	if err != nil {
		e.logger.Error("Failed to create output file: %v", err)
		return fmt.Errorf("failed to create file %s: %v", outputPath, err)
	}
	defer outFile.Close()
	
	// Copy with progress tracking
	e.logger.Info("Starting file transfer...")
	bytesWritten, err := io.Copy(outFile, resp.Body)
	if err != nil {
		e.logger.Error("File transfer failed: %v", err)
		return fmt.Errorf("failed to save download: %v", err)
	}
	
	e.logger.Info("Download completed successfully")
	e.logger.Info("Downloaded %d bytes to: %s", bytesWritten, outputPath)
	return nil
}

// handleMoveCommand implements file moving/renaming with verbose logging
func (e *Engine) handleMoveCommand(command, workDir string) error {
	e.logger.Info("Starting file move operation")
	
	// Parse MOVE: <source> to: <destination> or MOVE: <source> to <destination>
	parts := strings.SplitN(command, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid MOVE command format: %s", command)
	}
	
	moveArgs := strings.TrimSpace(parts[1])
	e.logger.Debug("Move arguments: %s", moveArgs)
	
	// Split by "to:" or "to "
	var source, dest string
	if strings.Contains(moveArgs, " to: ") {
		moveParts := strings.SplitN(moveArgs, " to: ", 2)
		source = strings.TrimSpace(moveParts[0])
		dest = strings.TrimSpace(moveParts[1])
	} else if strings.Contains(moveArgs, " to ") {
		moveParts := strings.SplitN(moveArgs, " to ", 2)
		source = strings.TrimSpace(moveParts[0])
		dest = strings.TrimSpace(moveParts[1])
	} else {
		return fmt.Errorf("invalid MOVE command format, missing 'to': %s", command)
	}
	
	if source == "" || dest == "" {
		return fmt.Errorf("empty source or destination in MOVE command")
	}
	
	e.logger.Info("Source: %s", source)
	e.logger.Info("Destination: %s", dest)
	
	// Make source path relative to workDir if not absolute
	if !filepath.IsAbs(source) {
		source = filepath.Join(workDir, source)
		e.logger.Debug("Resolved source path: %s", source)
	}
	
	// Check if source exists
	if _, err := os.Stat(source); os.IsNotExist(err) {
		e.logger.Error("Source file/directory does not exist: %s", source)
		return fmt.Errorf("source does not exist: %s", source)
	}
	
	// Make destination directory if needed
	destDir := filepath.Dir(dest)
	e.logger.Debug("Ensuring destination directory exists: %s", destDir)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		e.logger.Error("Failed to create destination directory: %v", err)
		return fmt.Errorf("failed to create destination directory %s: %v", destDir, err)
	}
	
	e.logger.Info("Moving: %s -> %s", source, dest)
	
	if err := os.Rename(source, dest); err != nil {
		e.logger.Error("Move operation failed: %v", err)
		return fmt.Errorf("failed to move %s to %s: %v", source, dest, err)
	}
	
	e.logger.Info("Move operation completed successfully")
	return nil
}

// handleRemoveCommand implements safe directory/file deletion with verbose logging
func (e *Engine) handleRemoveCommand(command, workDir string) error {
	e.logger.Info("Starting file/directory removal operation")
	
	// Parse REMOVE: <path>
	parts := strings.SplitN(command, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid REMOVE command format: %s", command)
	}
	
	targetPath := strings.TrimSpace(parts[1])
	if targetPath == "" {
		return fmt.Errorf("empty path in REMOVE command")
	}
	
	e.logger.Info("Target path: %s", targetPath)
	
	// Security checks
	if targetPath == "/" || targetPath == "C:\\" {
		e.logger.Error("Refusing to remove root directory: %s", targetPath)
		return fmt.Errorf("refusing to remove root directory")
	}
	
	// Expand path if not absolute
	if !filepath.IsAbs(targetPath) {
		targetPath = filepath.Join(workDir, targetPath)
		e.logger.Debug("Resolved target path: %s", targetPath)
	}
	
	// Check if path exists
	info, err := os.Stat(targetPath)
	if os.IsNotExist(err) {
		e.logger.Warn("Path does not exist, skipping removal: %s", targetPath)
		return nil
	}
	
	if info.IsDir() {
		e.logger.Info("Target is a directory")
	} else {
		e.logger.Info("Target is a file (size: %d bytes)", info.Size())
	}
	
	e.logger.Info("Removing: %s", targetPath)
	
	if err := os.RemoveAll(targetPath); err != nil {
		e.logger.Error("Remove operation failed: %v", err)
		return fmt.Errorf("failed to remove %s: %v", targetPath, err)
	}
	
	e.logger.Info("Remove operation completed successfully")
	return nil
} 