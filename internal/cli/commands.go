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
}
