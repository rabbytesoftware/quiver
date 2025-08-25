# Execution Engine

The execution engine is responsible for running arrow methods with support for special commands, hardware validation, and background execution.

## Features

### Variable System

The execution engine automatically provides system variables that arrows can reference:

#### INSTALL_DIR

The `INSTALL_DIR` variable is automatically created by Quiver for each arrow. It represents the directory where all files related to the arrow should be installed and managed.

- **Location**: Child directory of the `install_dir` configured in `config.json`
- **Naming**: Uses the arrow filename without extension (e.g., `chat.yaml` → `./pkgs/chat/`)
- **Usage**: Reference as `${INSTALL_DIR}` in arrow commands

Example:
```yaml
methods:
  windows:
    amd64:
      install:
        - GET: "https://example.com/app.exe"
        - MOVE: "app.exe" to: "${INSTALL_DIR}/app.exe"
      execute:
        - "${INSTALL_DIR}/app.exe --port ${APP_PORT}"
```

**Note**: For backward compatibility, `${INSTALL_PATH}` is also supported and points to the same location as `${INSTALL_DIR}`.

### Special Command Interpreters

The execution engine supports several special commands that are interpreted natively instead of being passed to the shell:

#### GET: URL

Downloads files from HTTP(S) URLs.

```yaml
- GET: "https://example.com/file.zip"
```

#### UNCOMPRESS: filename

Extracts compressed archives with support for multiple formats:

- **.zip** - Standard ZIP archives
- **.tar** - Uncompressed TAR archives
- **.tar.gz/.tgz** - Gzipped TAR archives
- **.rar** - WinRAR archives
- **.7z** - 7-Zip archives

```yaml
- UNCOMPRESS: "file.zip"
- UNCOMPRESS: "archive.tar.gz"
- UNCOMPRESS: "data.rar"
- UNCOMPRESS: "package.7z"
- UNCOMPRESS: "backup.tar"
```

All extractions include security validation to prevent directory traversal attacks.

#### MOVE: source to: destination

Moves or renames files and directories.

```yaml
- MOVE: "file.exe" to: "/usr/local/bin/app"
- MOVE: "source.txt" to "destination.txt"
```

#### REMOVE: path

Safely removes files and directories with security checks.

```yaml
- REMOVE: "/tmp/downloads"
- REMOVE: "unwanted-file.txt"
```

### Architecture Detection

The engine automatically detects the runtime OS and architecture (`runtime.GOOS` and `runtime.GOARCH`) and selects the appropriate method commands. If the exact architecture isn't available, it falls back to compatible alternatives.

### Hardware Requirement Validation

Before executing methods, the engine validates system requirements:

- **CPU Cores**: Checks available CPU cores using `gopsutil`
- **RAM**: Validates available memory in GB
- **Disk Space**: Ensures sufficient free disk space in GB

Requirements are defined in the arrow manifest's `requirements.minimum` section.

### Execution Options

The engine supports various execution modes:

- **Foreground**: Default mode, waits for command completion
- **Background**: Starts commands and returns immediately
- **Daemon**: For long-running services
- **Dry Run**: Simulates execution without running commands
- **Timeout**: Configurable timeouts for command execution

### Security Features

- **Path Validation**: Archive extraction prevents directory traversal attacks
- **Root Protection**: Refuses to remove system root directories
- **URL Validation**: HTTP requests include proper timeout and context handling
- **Environment Isolation**: Commands run with controlled environment variables
- **Archive Safety**: All archive formats are extracted with path sanitization

### Supported Archive Formats


| Format | Extension         | Library Used                      | Description               |
| ------ | ----------------- | --------------------------------- | ------------------------- |
| ZIP    | `.zip`            | `archive/zip`                     | Standard ZIP compression  |
| TAR    | `.tar`            | `archive/tar`                     | Uncompressed TAR archives |
| TAR.GZ | `.tar.gz`, `.tgz` | `archive/tar` + `compress/gzip`   | Gzipped TAR archives      |
| RAR    | `.rar`            | `github.com/nwaples/rardecode/v2` | WinRAR format             |
| 7-Zip  | `.7z`             | `github.com/bodgit/sevenzip`      | 7-Zip compression         |

## Usage

```go
// Create engine
engine := execution.NewEngine(logger)

// Basic execution
err := engine.ExecuteMethod(arrow, types.MethodInstall, ctx)

// Advanced execution with options
opts := &execution.ExecutionOptions{
    Mode:       execution.ExecutionModeBackground,
    Timeout:    10 * time.Minute,
    Background: true,
    DryRun:     false,
}
err := engine.ExecuteMethodWithOptions(arrow, types.MethodExecute, ctx, opts)

// Cleanup temporary files
defer engine.Cleanup()
```

## Error Handling

The engine provides detailed error messages for:

- Unsupported platforms/architectures
- Hardware requirement failures
- Command execution failures
- Network download issues
- File system operations
- Archive extraction errors
- Unsupported archive formats

All errors include context about what operation failed and why.

## Examples

### Complete Arrow Installation Flow

```yaml
methods:
  windows:
    amd64:
      install:
        - GET: "https://example.com/app-windows-amd64.zip"
        - UNCOMPRESS: "app-windows-amd64.zip" 
        - MOVE: "app.exe" to: "${INSTALL_DIR}\app.exe"
        - REMOVE: "app-windows-amd64.zip"
      execute:
        - "${INSTALL_DIR}\app.exe --port ${APP_PORT}"
```

This flow demonstrates downloading, extracting, moving, and cleaning up installation files using the special command interpreters.
