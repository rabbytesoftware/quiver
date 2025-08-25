# Netbridge Variables

## Overview

Netbridge variables are a powerful feature in Quiver that automatically manages port assignments for arrows. They bridge the gap between port forwarding capabilities and arrow execution by **dynamically assigning and opening ports at runtime** when arrows are executed.

## How It Works

### Basic Concept

Netbridge variables are defined in an arrow's manifest and represent ports that the system should attempt to open **during execution runtime**. The key principle is:

- **Installation**: No ports are assigned or opened
- **Initialization**: Variables are identified but no ports assigned  
- **Runtime Execution**: Ports are dynamically assigned and opened just before the arrow starts
- **Success case**: Port opens successfully → variable becomes available with the assigned port number
- **Failure case**: Port cannot be opened → variable still gets a fallback port number, but user is informed

This approach ensures that:
- Installation never opens network ports (security best practice)
- Ports are checked for availability at execution time (handles port conflicts)
- Multiple attempts can be made if initial ports are busy
- Arrows can always execute, with manual port configuration as fallback

### Workflow

1. **Arrow Definition**: Arrow manifest includes `netbridge` section with port specifications
2. **Installation**: Arrow is installed with no port assignment or opening
3. **Method Initialization**: When a method is initialized, netbridge variables are identified but no ports assigned
4. **Runtime Execution**: Just before executing arrow commands:
   - System attempts to find available ports
   - Attempts to open each netbridge port using UPnP/NAT-PMP  
   - Port numbers are assigned to variables regardless of opening success
5. **Variable Expansion**: Arrow commands use `${VARIABLE_NAME}` syntax to access dynamically assigned port numbers
6. **Status Reporting**: Execution logs report which ports opened successfully and which failed

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
- **Note**: Netbridge variables are only assigned during execute methods

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

**Response (Identification Only)**:
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
    "netbridge": {
      "variables": 1,
      "identifications": [
        {
          "variable_name": "CHAT_PORT",
          "protocol": "tcp/udp",
          "user_specified": true,
          "user_value": "8080"
        }
      ],
      "note": "Ports will be assigned dynamically at runtime during execution"
    },
    "status": "initialized"
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
    "identifications": [
      {
        "variable_name": "CHAT_PORT",
        "protocol": "tcp/udp",
        "user_specified": false,
        "user_value": ""
      }
    ],
    "note": "Ports are assigned dynamically at runtime, not during initialization",
    "status": "identified"
  }
}
```

## User Workflow Examples

### Scenario 1: Successful Dynamic Port Assignment

1. User calls `POST /api/v1/arrows/quiver.chat/initialize/execute`
2. System identifies netbridge variables but assigns no ports
3. User sees identification status in API response
4. User calls `POST /api/v1/arrows/quiver.chat/execute` to run the arrow
5. **At runtime**: System finds available port (e.g., 8080) and opens it via UPnP
6. `CHAT_PORT` variable is set to "8080" and expanded in commands
7. Arrow executes with dynamically assigned port

### Scenario 2: Port Opening Failure with Fallback

1. User calls `POST /api/v1/arrows/quiver.chat/initialize/execute`
2. System identifies netbridge variables (no ports assigned)
3. User calls `POST /api/v1/arrows/quiver.chat/execute`
4. **At runtime**: System tries to open port 8080 but UPnP/NAT-PMP fails
5. `CHAT_PORT` variable is still set to "8080" (fallback)
6. Arrow executes with assigned port (app must handle closed port)
7. Execution logs show port opening failure

### Scenario 3: User-Specified Port with Runtime Check

1. User calls `POST /api/v1/arrows/quiver.chat/initialize/execute` with `{"variables": {"CHAT_PORT": "9000"}}`
2. System identifies user wants port 9000 (no assignment yet)
3. User calls `POST /api/v1/arrows/quiver.chat/execute`
4. **At runtime**: System checks if port 9000 is available and tries to open it
5. `CHAT_PORT` variable is set to "9000" regardless of opening result
6. Arrow executes with user-specified port

## Port Assignment Logic

### Dynamic Assignment at Runtime

When no user port is specified:

1. **At runtime**: Try `netbridge.OpenPortAuto()` with specified protocol
2. If successful: use the auto-assigned port number
3. If failed: find an available port in range 8000-9000 without opening
4. If no available port: default to 8080

### User-Specified Ports

When user provides a port number:

1. **At runtime**: Validate port number (1-65535)
2. Try `netbridge.OpenPort()` with user's port and protocol
3. Use the user's port number regardless of opening success
4. Report opening status in execution logs

### Port Range

- Auto-assignment range: 8000-9000
- User can specify any valid port: 1-65535
- System avoids well-known ports (1-1023) in auto-assignment

## Error Handling

### Common Scenarios

1. **Netbridge Unavailable**: UPnP/NAT-PMP not supported on network
   - Variables still get assigned port numbers at runtime
   - Error logged during execution
   - Arrow can still execute

2. **Port Already in Use**: Specified port is occupied at runtime
   - Auto-assignment finds alternative port
   - User-specified ports report error but keep assignment

3. **Invalid Port Number**: User provides invalid port
   - Falls back to auto-assignment at runtime
   - Logs warning about invalid input

4. **Network Timeout**: UPnP/NAT-PMP requests timeout at runtime
   - Treated as opening failure
   - Variables still assigned
   - Error logged during execution

## Best Practices

### For Arrow Developers

1. **Always provide fallbacks**: Handle cases where ports can't be opened
2. **Use descriptive names**: `CHAT_PORT` vs `PORT1` 
3. **Document port usage**: Explain what each netbridge variable is for
4. **Test both scenarios**: Verify arrow works with opened and unopened ports
5. **Handle runtime failures**: Applications should gracefully handle port assignment failures

### For Users

1. **Check execution logs**: Review netbridge status during arrow execution
2. **Understand port requirements**: Know which ports are critical vs optional
3. **Manual configuration**: Be prepared to manually open failed ports if needed
4. **Security considerations**: Understand which ports are being opened dynamically

### For System Administrators

1. **Network compatibility**: Ensure UPnP/NAT-PMP is available if automatic opening is desired
2. **Firewall rules**: Consider firewall implications of dynamically opened ports
3. **Monitoring**: Watch for failed port openings in execution logs
4. **Documentation**: Inform users about network capabilities and limitations

## Security Considerations

### Dynamic Port Opening

- **Benefit**: No ports opened during installation (reduced attack surface)
- **Risk**: Ports opened automatically during execution
- **Mitigation**: User always informed via execution logs
- **Best Practice**: Review opened ports periodically and monitor execution logs

### Port Assignment

- **Benefit**: Ports checked for availability at runtime
- **Risk**: Predictable port numbers in auto-assignment range
- **Mitigation**: Use specific port assignment when security is critical
- **Best Practice**: Combine with additional security measures (authentication, encryption)

### Network Exposure

- **Benefit**: Ports only open when applications are actually running
- **Risk**: Services become externally accessible during execution
- **Mitigation**: Ensure applications implement proper security
- **Best Practice**: Use authentication, monitor access logs, and close ports when not needed

## Implementation Details

### Timing Overview

```
Installation Time:
├── Arrow files downloaded
├── Installation commands executed  
└── NO netbridge processing

Initialization Time:
├── Method validation
├── Variable identification (netbridge variables identified but no ports assigned)
└── API response with identification data

Runtime Execution:
├── Dynamic port assignment starts
├── Port availability check
├── UPnP/NAT-PMP port opening attempts
├── Variable assignment (regardless of opening success)
├── Command expansion with assigned port numbers
└── Arrow execution begins
```

### Error Recovery

The system is designed to always allow arrow execution:

1. **Primary**: Try to open requested/auto-assigned port
2. **Fallback 1**: Assign port number even if opening fails
3. **Fallback 2**: Use default port 8080 if no ports available
4. **Result**: Arrow always gets port numbers, execution proceeds

This ensures reliability while providing best-effort port opening. 