# Quiver CLI

The Quiver CLI allows you to interact with the Quiver server using command line arguments instead of making direct HTTP requests.

## Usage Modes

### Server Mode
When run without arguments, Quiver starts the REST API server and runs continuously:

```bash
./quiver
```

### CLI Mode
When run with arguments, Quiver starts the server temporarily, executes your command, and then exits:

```bash
./quiver [command] [arguments...]
```

## Available Commands

### Server Information
- `health` - Check server health status
- `info` - Get server information  
- `server-status` - Get server status

### Package Management
- `list` - List all packages
- `get <package_id>` - Get information about a specific package
- `start <package_id>` - Start a package
- `stop <package_id>` - Stop a package  
- `status <package_id>` - Get package status

### Help
- `help` - Show help information

## Examples

```bash
# Check server health
./quiver health

# List all packages
./quiver list

# Get information about a specific package
./quiver get minecraft-server

# Start a package
./quiver start minecraft-server

# Stop a package
./quiver stop minecraft-server

# Check package status
./quiver status minecraft-server

# Get server information
./quiver info

# Show help
./quiver help
```

## Command Implementation

Each CLI command is mapped to a REST API endpoint:

| CLI Command | HTTP Method | REST Endpoint |
|-------------|-------------|---------------|
| `health` | GET | `/health` |
| `list` | GET | `/api/v1/packages` |
| `get <id>` | GET | `/api/v1/packages/{id}` |
| `start <id>` | POST | `/api/v1/packages/{id}/start` |
| `stop <id>` | POST | `/api/v1/packages/{id}/stop` |
| `status <id>` | GET | `/api/v1/packages/{id}/status` |
| `info` | GET | `/api/v1/server/info` |
| `server-status` | GET | `/api/v1/server/status` |

## Adding New Commands

To add a new CLI command:

1. Add an entry to the `CommandRegistry` in `cli.go`
2. Define the command name, description, HTTP method, endpoint, and parameter types
3. The CLI will automatically handle parameter validation and HTTP request construction

Example:
```go
"new-command": {
    Name:        "new-command",
    Description: "Description of the new command",
    Method:      "GET",
    Endpoint:    "/api/v1/new-endpoint",
    ParamTypes: []ParamType{
        {Name: "param1", Type: "string", Required: true, Position: 0},
    },
    Example: "./quiver new-command param1-value",
},
``` 