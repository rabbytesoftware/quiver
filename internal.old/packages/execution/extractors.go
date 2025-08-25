package execution

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/nwaples/rardecode/v2"
)

// handleUncompressCommand implements extraction for multiple archive formats with verbose logging
func (e *Engine) handleUncompressCommand(command, workDir string) error {
	e.logger.Info("Starting archive extraction operation")
	
	// Parse UNCOMPRESS: <filename>
	parts := strings.SplitN(command, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid UNCOMPRESS command format: %s", command)
	}
	
	filename := strings.TrimSpace(parts[1])
	if filename == "" {
		return fmt.Errorf("empty filename in UNCOMPRESS command")
	}
	
	archivePath := filepath.Join(workDir, filename)
	
	e.logger.Info("Archive file: %s", filename)
	e.logger.Info("Full path: %s", archivePath)
	e.logger.Info("Working directory: %s", workDir)
	
	// Check if file exists and get file info
	fileInfo, err := os.Stat(archivePath)
	if os.IsNotExist(err) {
		e.logger.Error("Archive file does not exist: %s", archivePath)
		return fmt.Errorf("archive file does not exist: %s", archivePath)
	}
	
	e.logger.Info("Archive size: %d bytes", fileInfo.Size())
	
	// Determine format by file extension
	ext := strings.ToLower(filepath.Ext(filename))
	if strings.HasSuffix(strings.ToLower(filename), ".tar.gz") {
		ext = ".tar.gz"
	}
	
	e.logger.Info("Detected archive format: %s", ext)
	
	var fileCount int
	
	switch ext {
	case ".zip":
		e.logger.Info("Using ZIP extractor")
		fileCount, err = e.extractZip(archivePath, workDir)
	case ".tar":
		e.logger.Info("Using TAR extractor")
		fileCount, err = e.extractTar(archivePath, workDir)
	case ".tar.gz", ".tgz":
		e.logger.Info("Using TAR.GZ extractor")
		fileCount, err = e.extractTarGz(archivePath, workDir)
	case ".rar":
		e.logger.Info("Using RAR extractor")
		fileCount, err = e.extractRar(archivePath, workDir)
	case ".7z":
		e.logger.Info("Using 7-Zip extractor")
		fileCount, err = e.extract7z(archivePath, workDir)
	default:
		e.logger.Error("Unsupported archive format: %s", ext)
		return fmt.Errorf("unsupported archive format: %s (supported: .zip, .tar, .tar.gz, .tgz, .rar, .7z)", ext)
	}
	
	if err != nil {
		e.logger.Error("Extraction failed: %v", err)
		return fmt.Errorf("failed to extract %s: %v", filename, err)
	}
	
	e.logger.Info("Archive extraction completed successfully")
	e.logger.Info("Extracted %d files from %s", fileCount, filename)
	return nil
}

// extractZip extracts a ZIP archive with verbose logging
func (e *Engine) extractZip(archivePath, workDir string) (int, error) {
	e.logger.Debug("Opening ZIP archive: %s", archivePath)
	
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		e.logger.Error("Failed to open ZIP file: %v", err)
		return 0, fmt.Errorf("failed to open ZIP file: %v", err)
	}
	defer reader.Close()
	
	totalFiles := len(reader.File)
	e.logger.Info("ZIP archive contains %d files", totalFiles)
	
	for i, file := range reader.File {
		e.logger.Debug("Extracting file %d/%d: %s", i+1, totalFiles, file.Name)
		
		if err := e.extractZipFile(file, workDir); err != nil {
			e.logger.Error("Failed to extract file %s: %v", file.Name, err)
			return 0, fmt.Errorf("failed to extract %s: %v", file.Name, err)
		}
	}
	
	e.logger.Info("ZIP extraction completed: %d files extracted", totalFiles)
	return totalFiles, nil
}

// extractZipFile extracts a single file from a ZIP archive
func (e *Engine) extractZipFile(file *zip.File, workDir string) error {
	path := filepath.Join(workDir, file.Name)
	
	// Security check: ensure path is within workDir
	if !isPathSafe(path, workDir) {
		e.logger.Error("Security violation: path outside working directory: %s", file.Name)
		return fmt.Errorf("invalid file path in archive: %s", file.Name)
	}
	
	if file.FileInfo().IsDir() {
		e.logger.Debug("Creating directory: %s", path)
		return os.MkdirAll(path, 0755)
	}
	
	// Create directory for file
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	
	outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()
	
	bytesWritten, err := io.Copy(outFile, rc)
	if err != nil {
		return err
	}
	
	e.logger.Debug("Extracted file: %s (%d bytes)", file.Name, bytesWritten)
	return nil
}

// extractTar extracts a TAR archive with verbose logging
func (e *Engine) extractTar(archivePath, workDir string) (int, error) {
	e.logger.Debug("Opening TAR archive: %s", archivePath)
	
	file, err := os.Open(archivePath)
	if err != nil {
		e.logger.Error("Failed to open TAR file: %v", err)
		return 0, fmt.Errorf("failed to open TAR file: %v", err)
	}
	defer file.Close()
	
	tarReader := tar.NewReader(file)
	return e.extractTarReader(tarReader, workDir)
}

// extractTarGz extracts a gzipped TAR archive with verbose logging
func (e *Engine) extractTarGz(archivePath, workDir string) (int, error) {
	e.logger.Debug("Opening TAR.GZ archive: %s", archivePath)
	
	file, err := os.Open(archivePath)
	if err != nil {
		e.logger.Error("Failed to open TAR.GZ file: %v", err)
		return 0, fmt.Errorf("failed to open TAR.GZ file: %v", err)
	}
	defer file.Close()
	
	e.logger.Debug("Creating gzip reader")
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		e.logger.Error("Failed to create gzip reader: %v", err)
		return 0, fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()
	
	tarReader := tar.NewReader(gzReader)
	return e.extractTarReader(tarReader, workDir)
}

// extractTarReader extracts files from a TAR reader with verbose logging
func (e *Engine) extractTarReader(tarReader *tar.Reader, workDir string) (int, error) {
	fileCount := 0
	
	e.logger.Debug("Starting TAR extraction")
	
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			e.logger.Error("Failed to read TAR header: %v", err)
			return 0, fmt.Errorf("failed to read TAR header: %v", err)
		}
		
		fileCount++
		path := filepath.Join(workDir, header.Name)
		
		e.logger.Debug("Processing TAR entry %d: %s", fileCount, header.Name)
		
		// Security check: ensure path is within workDir
		if !isPathSafe(path, workDir) {
			e.logger.Error("Security violation: path outside working directory: %s", header.Name)
			return 0, fmt.Errorf("invalid file path in archive: %s", header.Name)
		}
		
		switch header.Typeflag {
		case tar.TypeDir:
			e.logger.Debug("Creating directory: %s", path)
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				e.logger.Error("Failed to create directory %s: %v", path, err)
				return 0, fmt.Errorf("failed to create directory %s: %v", path, err)
			}
		case tar.TypeReg:
			// Create directory for file
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return 0, fmt.Errorf("failed to create directory: %v", err)
			}
			
			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				e.logger.Error("Failed to create file %s: %v", path, err)
				return 0, fmt.Errorf("failed to create file %s: %v", path, err)
			}
			
			bytesWritten, err := io.Copy(outFile, tarReader)
			outFile.Close()
			
			if err != nil {
				e.logger.Error("Failed to extract file %s: %v", path, err)
				return 0, fmt.Errorf("failed to extract file %s: %v", path, err)
			}
			
			e.logger.Debug("Extracted file: %s (%d bytes)", header.Name, bytesWritten)
		default:
			e.logger.Warn("Skipping unsupported file type %c: %s", header.Typeflag, header.Name)
		}
	}
	
	e.logger.Info("TAR extraction completed: %d entries processed", fileCount)
	return fileCount, nil
}

// extractRar extracts a RAR archive with verbose logging
func (e *Engine) extractRar(archivePath, workDir string) (int, error) {
	e.logger.Debug("Opening RAR archive: %s", archivePath)
	
	file, err := os.Open(archivePath)
	if err != nil {
		e.logger.Error("Failed to open RAR file: %v", err)
		return 0, fmt.Errorf("failed to open RAR file: %v", err)
	}
	defer file.Close()
	
	reader, err := rardecode.NewReader(file)
	if err != nil {
		e.logger.Error("Failed to create RAR reader: %v", err)
		return 0, fmt.Errorf("failed to create RAR reader: %v", err)
	}
	
	fileCount := 0
	
	e.logger.Debug("Starting RAR extraction")
	
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			e.logger.Error("Failed to read RAR header: %v", err)
			return 0, fmt.Errorf("failed to read RAR header: %v", err)
		}
		
		fileCount++
		path := filepath.Join(workDir, header.Name)
		
		e.logger.Debug("Processing RAR entry %d: %s", fileCount, header.Name)
		
		// Security check: ensure path is within workDir
		if !isPathSafe(path, workDir) {
			e.logger.Error("Security violation: path outside working directory: %s", header.Name)
			return 0, fmt.Errorf("invalid file path in archive: %s", header.Name)
		}
		
		if header.IsDir {
			e.logger.Debug("Creating directory: %s", path)
			if err := os.MkdirAll(path, 0755); err != nil {
				e.logger.Error("Failed to create directory %s: %v", path, err)
				return 0, fmt.Errorf("failed to create directory %s: %v", path, err)
			}
		} else {
			// Create directory for file
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return 0, fmt.Errorf("failed to create directory: %v", err)
			}
			
			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, header.Mode())
			if err != nil {
				e.logger.Error("Failed to create file %s: %v", path, err)
				return 0, fmt.Errorf("failed to create file %s: %v", path, err)
			}
			
			bytesWritten, err := io.Copy(outFile, reader)
			outFile.Close()
			
			if err != nil {
				e.logger.Error("Failed to extract file %s: %v", path, err)
				return 0, fmt.Errorf("failed to extract file %s: %v", path, err)
			}
			
			e.logger.Debug("Extracted file: %s (%d bytes)", header.Name, bytesWritten)
		}
	}
	
	e.logger.Info("RAR extraction completed: %d entries processed", fileCount)
	return fileCount, nil
}

// extract7z extracts a 7z archive with verbose logging
func (e *Engine) extract7z(archivePath, workDir string) (int, error) {
	e.logger.Debug("Opening 7-Zip archive: %s", archivePath)
	
	reader, err := sevenzip.OpenReader(archivePath)
	if err != nil {
		e.logger.Error("Failed to open 7z file: %v", err)
		return 0, fmt.Errorf("failed to open 7z file: %v", err)
	}
	defer reader.Close()
	
	totalFiles := len(reader.File)
	e.logger.Info("7-Zip archive contains %d files", totalFiles)
	
	for i, file := range reader.File {
		e.logger.Debug("Extracting file %d/%d: %s", i+1, totalFiles, file.Name)
		
		if err := e.extract7zFile(file, workDir); err != nil {
			e.logger.Error("Failed to extract file %s: %v", file.Name, err)
			return 0, fmt.Errorf("failed to extract %s: %v", file.Name, err)
		}
	}
	
	e.logger.Info("7-Zip extraction completed: %d files extracted", totalFiles)
	return totalFiles, nil
}

// extract7zFile extracts a single file from a 7z archive
func (e *Engine) extract7zFile(file *sevenzip.File, workDir string) error {
	path := filepath.Join(workDir, file.Name)
	
	// Security check: ensure path is within workDir
	if !isPathSafe(path, workDir) {
		e.logger.Error("Security violation: path outside working directory: %s", file.Name)
		return fmt.Errorf("invalid file path in archive: %s", file.Name)
	}
	
	if file.FileInfo().IsDir() {
		e.logger.Debug("Creating directory: %s", path)
		return os.MkdirAll(path, 0755)
	}
	
	// Create directory for file
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	
	outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()
	
	bytesWritten, err := io.Copy(outFile, rc)
	if err != nil {
		return err
	}
	
	e.logger.Debug("Extracted file: %s (%d bytes)", file.Name, bytesWritten)
	return nil
} 