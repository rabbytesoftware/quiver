package cli

// CommandRegistry holds all available commands
var CommandRegistry = map[string]Command{
	"health": {
		Name:        "health",
		Description: "Check server health status",
		Method:      "GET",
		Endpoint:    "/health",
		ParamTypes:  []ParamType{},
		Example:     "./quiver health",
	},
	"list": {
		Name:        "list",
		Description: "List all packages",
		Method:      "GET", 
		Endpoint:    "/api/v1/packages",
		ParamTypes:  []ParamType{},
		Example:     "./quiver list",
	},
	"get": {
		Name:        "get",
		Description: "Get information about a specific package",
		Method:      "GET",
		Endpoint:    "/api/v1/packages/{id}",
		ParamTypes: []ParamType{
			{Name: "package_id", Type: "string", Required: true, Position: 0, URLParam: true},
		},
		Example: "./quiver get minecraft-server",
	},
	"start": {
		Name:        "start",
		Description: "Start a package",
		Method:      "POST",
		Endpoint:    "/api/v1/packages/{id}/start",
		ParamTypes: []ParamType{
			{Name: "package_id", Type: "string", Required: true, Position: 0, URLParam: true},
		},
		Example: "./quiver start minecraft-server",
	},
	"stop": {
		Name:        "stop",
		Description: "Stop a package",
		Method:      "POST",
		Endpoint:    "/api/v1/packages/{id}/stop",
		ParamTypes: []ParamType{
			{Name: "package_id", Type: "string", Required: true, Position: 0, URLParam: true},
		},
		Example: "./quiver stop minecraft-server",
	},
	"status": {
		Name:        "status",
		Description: "Get package status",
		Method:      "GET",
		Endpoint:    "/api/v1/packages/{id}/status",
		ParamTypes: []ParamType{
			{Name: "package_id", Type: "string", Required: true, Position: 0, URLParam: true},
		},
		Example: "./quiver status minecraft-server",
	},
	"info": {
		Name:        "info",
		Description: "Get server information",
		Method:      "GET",
		Endpoint:    "/api/v1/server/info",
		ParamTypes:  []ParamType{},
		Example:     "./quiver info",
	},
	"server-status": {
		Name:        "server-status",
		Description: "Get server status",
		Method:      "GET",
		Endpoint:    "/api/v1/server/status",
		ParamTypes:  []ParamType{},
		Example:     "./quiver server-status",
	},
	// Arrow package management commands
	"search": {
		Name:        "search",
		Description: "Search for arrows in repositories",
		Method:      "GET",
		Endpoint:    "/api/v1/arrows/search",
		ParamTypes: []ParamType{
			{Name: "query", Type: "string", Required: false, Position: 0, URLParam: false},
		},
		Example: "./quiver search minecraft",
	},
	"install": {
		Name:        "install",
		Description: "Install an arrow",
		Method:      "POST",
		Endpoint:    "/api/v1/arrows/{name}/install",
		ParamTypes: []ParamType{
			{Name: "name", Type: "string", Required: true, Position: 0, URLParam: true},
		},
		Example: "./quiver install cs2",
	},
	"execute": {
		Name:        "execute",
		Description: "Execute an installed arrow",
		Method:      "POST",
		Endpoint:    "/api/v1/arrows/{name}/execute",
		ParamTypes: []ParamType{
			{Name: "name", Type: "string", Required: true, Position: 0, URLParam: true},
		},
		Example: "./quiver execute cs2",
	},
	"uninstall": {
		Name:        "uninstall",
		Description: "Uninstall an arrow",
		Method:      "DELETE",
		Endpoint:    "/api/v1/arrows/{name}/uninstall",
		ParamTypes: []ParamType{
			{Name: "name", Type: "string", Required: true, Position: 0, URLParam: true},
		},
		Example: "./quiver uninstall cs2",
	},
	"update": {
		Name:        "update",
		Description: "Update an arrow to the latest version",
		Method:      "PUT",
		Endpoint:    "/api/v1/arrows/{name}/update",
		ParamTypes: []ParamType{
			{Name: "name", Type: "string", Required: true, Position: 0, URLParam: true},
		},
		Example: "./quiver update cs2",
	},
	"validate": {
		Name:        "validate",
		Description: "Validate an arrow installation",
		Method:      "POST",
		Endpoint:    "/api/v1/arrows/{name}/validate",
		ParamTypes: []ParamType{
			{Name: "name", Type: "string", Required: true, Position: 0, URLParam: true},
		},
		Example: "./quiver validate cs2",
	},
	"installed": {
		Name:        "installed",
		Description: "List all installed arrows",
		Method:      "GET",
		Endpoint:    "/api/v1/arrows/installed",
		ParamTypes:  []ParamType{},
		Example:     "./quiver installed",
	},
	"arrow-status": {
		Name:        "arrow-status",
		Description: "Get arrow installation status",
		Method:      "GET",
		Endpoint:    "/api/v1/arrows/{name}/status",
		ParamTypes: []ParamType{
			{Name: "name", Type: "string", Required: true, Position: 0, URLParam: true},
		},
		Example: "./quiver arrow-status cs2",
	},
	// Repository management commands
	"repo-list": {
		Name:        "repo-list",
		Description: "List all repositories",
		Method:      "GET",
		Endpoint:    "/api/v1/repositories",
		ParamTypes:  []ParamType{},
		Example:     "./quiver repo-list",
	},
	"repo-add": {
		Name:        "repo-add",
		Description: "Add a new repository",
		Method:      "POST",
		Endpoint:    "/api/v1/repositories",
		ParamTypes: []ParamType{
			{Name: "repository", Type: "string", Required: true, Position: 0, URLParam: false},
		},
		Example: "./quiver repo-add https://github.com/example/arrows",
	},
	"repo-remove": {
		Name:        "repo-remove",
		Description: "Remove a repository",
		Method:      "DELETE",
		Endpoint:    "/api/v1/repositories",
		ParamTypes: []ParamType{
			{Name: "repository", Type: "string", Required: true, Position: 0, URLParam: false},
		},
		Example: "./quiver repo-remove https://github.com/example/arrows",
	},
}
