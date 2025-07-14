# Netbridge Variables

## Overview

Netbridge variables are a powerful feature in Quiver that automatically manages port assignments for arrows. They bridge the gap between port forwarding capabilities and arrow execution by attempting to open ports and making them available as variables to arrow methods.

## How It Works

### Basic Concept

Netbridge variables are defined in an arrow's manifest and represent ports that the system should attempt to open. The key principle is:

- **Success case**: Port opens successfully → variable becomes available with the assigned port number
- **Failure case**: Port cannot be opened → variable still gets the attempted port number, but user is informed via REST API

This approach ensures that arrows can always execute, while giving users control over manual port configuration when automatic opening fails.

### Workflow

1. **Arrow Definition**: Arrow manifest includes `netbridge` section with port specifications
2. **Method Initialization**: When a method is initialized, netbridge processing occurs
3. **Port Assignment**: System attempts to open each netbridge port using UPnP/NAT-PMP
4. **Variable Creation**: Port numbers are assigned to variables regardless of opening success
5. **Status Reporting**: REST API reports which ports opened successfully and which failed
6. **User Decision**: User can manually open failed ports or proceed with current assignment
7. **Method Execution**: Arrow methods use `${VARIABLE_NAME}` syntax to access port numbers

## Arrow Manifest Configuration

### Basic Syntax

```yaml
netbridge:
  - name: "CHAT_PORT"
    protocol: "tcp"
  - name: "API_PORT" 
    protocol: "tcp/udp"
  - name: "STREAMING_PORT"
    protocol: "udp"
```

### Supported Protocols

- `tcp`: TCP protocol only
- `udp`: UDP protocol only  
- `tcp/udp`: Both TCP and UDP (will attempt to open both)

### Variable Naming

- Use UPPERCASE with underscores (e.g., `CHAT_PORT`, `API_ENDPOINT_PORT`)
- Must be valid environment variable names
- Should be descriptive of the service/purpose

## Method Usage

### In Arrow Methods

```yaml
methods:
  execute:
    linux:
      - "${INSTALL_DIR}/app --port ${CHAT_PORT} --api-port ${API_PORT}"
    windows:
      - "${INSTALL_DIR}\\app.exe --port ${CHAT_PORT} --api-port ${API_PORT}"
```

### Variable Expansion

- Variables are expanded during command execution
- Use `${VARIABLE_NAME}` syntax
- Works in any arrow method (install, execute, uninstall, etc.)

## REST API Integration

### Method Initialization Endpoint

```http
POST /api/v1/arrows/:name/initialize/:method
Content-Type: application/json

{
  "variables": {
    "CHAT_PORT": "8080"
  }
}
```

**Response (Success)**:
```json
{
  "success": true,
  "message": "Method initialized successfully",
  "data": {
    "arrow": "quiver.chat",
    "method": "execute", 
    "variables": {
      "CHAT_PORT": "8080",
      "QUIVER_CHAT_HOSTNAME": "chat.quiver.ar"
    },
    "netbridge": [
      {
        "variable_name": "CHAT_PORT",
        "port": 8080,
        "protocol": "tcp/udp",
        "success": true,
        "result": {
          "port": {...},
          "method": "upnp",
          "success": true
        }
      }
    ],
    "status": "initialized"
  }
}
```

**Response (Partial Failure)**:
```json
{
  "success": true,
  "message": "Method initialized successfully",
  "data": {
    "arrow": "quiver.chat",
    "method": "execute",
    "variables": {
      "CHAT_PORT": "8080"
    },
    "netbridge": [
      {
        "variable_name": "CHAT_PORT", 
        "port": 8080,
        "protocol": "tcp/udp",
        "success": false,
        "error": "Port forwarding failed: UPnP and NAT-PMP methods unsuccessful"
      }
    ],
    "status": "initialized",
    "warnings": [
      "Port 8080 for variable CHAT_PORT could not be opened: Port forwarding failed: UPnP and NAT-PMP methods unsuccessful"
    ]
  }
}
```

### Netbridge Status Endpoint

```http
GET /api/v1/arrows/:name/netbridge
```

**Response**:
```json
{
  "success": true,
  "message": "Netbridge status retrieved successfully",
  "data": {
    "arrow": "quiver.chat",
    "netbridge_vars": 1,
    "variables": [
      {
        "variable_name": "CHAT_PORT",
        "port": 8080,
        "protocol": "tcp/udp", 
        "success": false,
        "error": "Netbridge not available"
      }
    ],
    "status": "processed"
  }
}
```

## User Workflow Examples

### Scenario 1: Successful Port Opening

1. User calls `POST /api/v1/arrows/quiver.chat/initialize/execute`
2. System successfully opens port 8080 via UPnP
3. `CHAT_PORT` variable is set to "8080" 
4. User sees success status in API response
5. User calls `POST /api/v1/arrows/quiver.chat/execute` to run the arrow
6. Arrow executes with `${CHAT_PORT}` expanded to "8080"

### Scenario 2: Failed Port Opening

1. User calls `POST /api/v1/arrows/quiver.chat/initialize/execute`
2. System fails to open port 8080 (UPnP/NAT-PMP unavailable)
3. `CHAT_PORT` variable is still set to "8080"
4. User sees failure status and warning in API response
5. **User Choice A**: Manually open port 8080 in router, then execute arrow
6. **User Choice B**: Accept that port is closed and execute anyway (app handles closed port)

### Scenario 3: User-Specified Port

1. User calls `POST /api/v1/arrows/quiver.chat/initialize/execute` with `{"variables": {"CHAT_PORT": "9000"}}`
2. System attempts to open user-specified port 9000
3. `CHAT_PORT` variable is set to "9000" regardless of opening result
4. User gets status about whether port 9000 was opened successfully

## Port Assignment Logic

### Auto-Assignment

When no user port is specified:

1. Try `netbridge.OpenPortAuto()` with specified protocol
2. If successful: use the auto-assigned port number
3. If failed: find an available port in range 8000-9000 without opening
4. If no available port: default to 8080

### User-Specified Ports

When user provides a port number:

1. Validate port number (1-65535)
2. Try `netbridge.OpenPort()` with user's port and protocol
3. Use the user's port number regardless of opening success
4. Report opening status to user

### Port Range

- Auto-assignment range: 8000-9000
- User can specify any valid port: 1-65535
- System avoids well-known ports (1-1023) in auto-assignment

## Error Handling

### Common Scenarios

1. **Netbridge Unavailable**: UPnP/NAT-PMP not supported on network
   - Variables still get assigned port numbers
   - Error reported to user
   - Arrow can still execute

2. **Port Already in Use**: Specified port is occupied
   - Auto-assignment finds alternative port
   - User-specified ports report error but keep assignment

3. **Invalid Port Number**: User provides invalid port
   - Falls back to auto-assignment
   - Logs warning about invalid input

4. **Network Timeout**: UPnP/NAT-PMP requests timeout
   - Treated as opening failure
   - Variables still assigned
   - Error reported to user

## Best Practices

### For Arrow Developers

1. **Always provide fallbacks**: Handle cases where ports can't be opened
2. **Use descriptive names**: `CHAT_PORT` vs `PORT1` 
3. **Document port usage**: Explain what each netbridge variable is for
4. **Test both scenarios**: Verify arrow works with opened and unopened ports

### For Users

1. **Check initialization results**: Review netbridge status before executing
2. **Understand port requirements**: Know which ports are critical vs optional
3. **Manual configuration**: Be prepared to manually open failed ports if needed
4. **Security considerations**: Understand which ports are being opened

### For System Administrators

1. **Network compatibility**: Ensure UPnP/NAT-PMP is available if automatic opening is desired
2. **Firewall rules**: Consider firewall implications of auto-opened ports
3. **Monitoring**: Watch for failed port openings in logs
4. **Documentation**: Inform users about network capabilities and limitations

## Integration Examples

### Complete Arrow Example

```yaml
version: "0.1"

metadata:
  name: "chat-server"
  description: "A simple chat server with API"
  version: "1.0.0"

netbridge:
  - name: "CHAT_PORT"
    protocol: "tcp"
  - name: "API_PORT"
    protocol: "tcp"
  - name: "STREAM_PORT"
    protocol: "udp"

variables:
  - name: "SERVER_NAME"
    default: "My Chat Server"

methods:
  execute:
    linux:
      - "${INSTALL_DIR}/chat-server --name '${SERVER_NAME}' --chat-port ${CHAT_PORT} --api-port ${API_PORT} --stream-port ${STREAM_PORT}"
```

### Client Application Usage

```javascript
// Initialize method with netbridge processing
const initResponse = await fetch('/api/v1/arrows/chat-server/initialize/execute', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    variables: {
      SERVER_NAME: "Production Chat"
    }
  })
});

const initData = await initResponse.json();

// Check netbridge results
for (const netVar of initData.data.netbridge) {
  if (!netVar.success) {
    console.warn(`Port ${netVar.port} for ${netVar.variable_name} failed to open: ${netVar.error}`);
    
    // Prompt user for manual configuration
    const shouldContinue = confirm(
      `Port ${netVar.port} could not be opened automatically. ` +
      `You may need to manually configure your router. Continue anyway?`
    );
    
    if (!shouldContinue) {
      return; // Don't execute the arrow
    }
  }
}

// Execute the arrow
const execResponse = await fetch('/api/v1/arrows/chat-server/execute', {
  method: 'POST'
});
```

## Troubleshooting

### Port Opening Failures

**Problem**: All netbridge variables show `success: false`

**Solutions**:
1. Check if UPnP is enabled on router
2. Verify NAT-PMP support
3. Check network connectivity
4. Review firewall settings
5. Consider manual port forwarding

**Problem**: Specific ports fail while others succeed

**Solutions**:
1. Check if port is already in use
2. Verify port is not blocked by ISP
3. Try different port numbers
4. Check port-specific router settings

### Variable Expansion Issues

**Problem**: `${VARIABLE_NAME}` not expanded in commands

**Solutions**:
1. Verify variable name matches netbridge definition
2. Check for typos in variable syntax
3. Ensure netbridge processing completed successfully
4. Review execution logs for variable assignment

### API Response Issues

**Problem**: Missing netbridge data in API responses

**Solutions**:
1. Ensure arrow has netbridge variables defined
2. Verify correct API endpoint usage
3. Check for arrow installation issues
4. Review server logs for processing errors

## Security Considerations

### Automatic Port Opening

- **Risk**: Opens firewall ports automatically
- **Mitigation**: User always informed of opening attempts
- **Best Practice**: Review opened ports periodically

### Port Assignment

- **Risk**: Predictable port numbers in auto-assignment
- **Mitigation**: Use specific port assignment when security is critical
- **Best Practice**: Combine with additional security measures (authentication, encryption)

### Network Exposure

- **Risk**: Services become externally accessible
- **Mitigation**: Ensure applications implement proper security
- **Best Practice**: Use authentication and monitor access logs 