# Quiver REST API Documentation

## Overview

The Quiver REST API provides comprehensive management capabilities for arrows, packages, repositories, and server monitoring. Built on the Gin framework with standardized responses and modular architecture.

## Base URL

```
http://localhost:8080
```

## Authentication

Currently, the API does not require authentication. Future versions will include API key or token-based authentication.

## Response Format

All API responses follow a standardized format:

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": "Additional error details (optional)"
  }
}
```

## HTTP Status Codes

- `200 OK` - Successful request
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request parameters
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Validation errors
- `500 Internal Server Error` - Server error

## API Endpoints

### Health Check Endpoints

#### Health Check
```
GET /health
```

**Response:**
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "service": "Quiver",
    "version": "1.0.0",
    "status": "healthy",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

#### Readiness Probe
```
GET /ready
```

#### Liveness Probe
```
GET /live
```

### Arrow Management

#### Search Arrows
```
GET /api/v1/arrows/search?q={query}
```

**Parameters:**
- `q` (required): Search query string

**Response:**
```json
{
  "success": true,
  "message": "Search completed successfully",
  "data": {
    "query": "minecraft",
    "results": [
      {
        "name": "minecraft-server",
        "description": "Minecraft server package",
        "version": "1.0.0",
        "repository": "https://github.com/example/arrows"
      }
    ]
  }
}
```

#### Install Arrow
```
POST /api/v1/arrows/{name}/install
```

**Request Body:**
```json
{
  "variables": {
    "SERVER_NAME": "My Server",
    "MAX_PLAYERS": 20
  }
}
```

**Response:**
```json
{
  "success": true,
  "message": "Arrow installed successfully",
  "data": {
    "name": "minecraft-server",
    "version": "1.0.0",
    "status": "installed",
    "installedAt": "2024-01-01T12:00:00Z"
  }
}
```

#### Execute Arrow
```
POST /api/v1/arrows/{name}/execute
```

#### Uninstall Arrow
```
DELETE /api/v1/arrows/{name}/uninstall
```

#### Update Arrow
```
PUT /api/v1/arrows/{name}/update
```

#### Validate Arrow
```
POST /api/v1/arrows/{name}/validate
```

#### Get Installed Arrows
```
GET /api/v1/arrows/installed
```

#### Get Arrow Status
```
GET /api/v1/arrows/{name}/status
```

#### Get Arrows by Status
```
GET /api/v1/arrows/status/{status}
```

**Parameters:**
- `status`: `running`, `stopped`, `installed`, `error`

#### Get Arrow Information
```
GET /api/v1/arrows/{name}
```

#### List Arrow Statuses
```
GET /api/v1/arrows/statuses
```

### Package Management

#### List Packages
```
GET /api/v1/packages/
```

**Response:**
```json
{
  "success": true,
  "message": "Available packages retrieved successfully",
  "data": {
    "packages": [
      {
        "id": "minecraft-server",
        "name": "Minecraft Server",
        "version": "1.0.0",
        "status": "installed"
      }
    ]
  }
}
```

#### Get Package Information
```
GET /api/v1/packages/{id}
```

#### Start Package
```
POST /api/v1/packages/{id}/start
```

#### Stop Package
```
POST /api/v1/packages/{id}/stop
```

#### Get Package Status
```
GET /api/v1/packages/{id}/status
```

**Response:**
```json
{
  "success": true,
  "message": "Package status retrieved successfully",
  "data": {
    "id": "minecraft-server",
    "status": "running",
    "uptime": "2h30m",
    "pid": 1234,
    "lastUpdate": "2024-01-01T12:00:00Z"
  }
}
```

#### Get Installed Packages
```
GET /api/v1/packages/installed
```

#### Get Packages by Status
```
GET /api/v1/packages/status/{status}
```

#### List Package Statuses
```
GET /api/v1/packages/statuses
```

### Repository Management

#### Add Repository
```
POST /api/v1/repositories/
```

**Request Body:**
```json
{
  "repository": "https://github.com/example/arrows",
  "branch": "main"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Repository added successfully",
  "data": {
    "repository": "https://github.com/example/arrows",
    "branch": "main",
    "addedAt": "2024-01-01T12:00:00Z"
  }
}
```

#### Remove Repository
```
DELETE /api/v1/repositories/
```

**Request Body:**
```json
{
  "repository": "https://github.com/example/arrows"
}
```

#### List Repositories
```
GET /api/v1/repositories/
```

**Response:**
```json
{
  "success": true,
  "message": "Repositories retrieved successfully",
  "data": {
    "repositories": [
      {
        "url": "https://github.com/example/arrows",
        "branch": "main",
        "lastUpdate": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

#### Get Repository Information
```
GET /api/v1/repositories/{repository}
```

#### Search Repositories
```
GET /api/v1/repositories/search?q={query}
```

#### Validate Repository
```
POST /api/v1/repositories/{repository}/validate
```

### Server Management

#### Get Server Information
```
GET /api/v1/server/info
```

**Response:**
```json
{
  "success": true,
  "message": "Server information retrieved successfully",
  "data": {
    "name": "Quiver Server",
    "version": "1.0.0",
    "buildTime": "2024-01-01T12:00:00Z",
    "goVersion": "go1.24.2",
    "platform": "linux/amd64"
  }
}
```

#### Get Server Status
```
GET /api/v1/server/status
```

#### Get Server Health
```
GET /api/v1/server/health
```

#### Get Server Metrics
```
GET /api/v1/server/metrics
```

**Response:**
```json
{
  "success": true,
  "message": "Server metrics retrieved successfully",
  "data": {
    "uptime": "24h30m",
    "requestsTotal": 1543,
    "requestsPerSecond": 12.5,
    "memoryUsage": "45MB",
    "cpuUsage": "2.3%",
    "goroutines": 42
  }
}
```

#### Get Server Version
```
GET /api/v1/server/version
```

### Network Bridge Management

#### Open Port
```
POST /api/v1/netbridge/open
```

**Request Body:**
```json
{
  "port": 25565,
  "protocol": "tcp",
  "description": "Minecraft server port"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Port opened successfully",
  "data": {
    "port": 25565,
    "protocol": "tcp",
    "publicIP": "203.0.113.1",
    "status": "open"
  }
}
```

#### Close Port
```
POST /api/v1/netbridge/close
```

**Request Body:**
```json
{
  "port": 25565,
  "protocol": "tcp"
}
```

#### Open Port by URL
```
POST /api/v1/netbridge/port/{port}/open
```

#### Close Port by URL
```
DELETE /api/v1/netbridge/port/{port}
```

#### List Open Ports
```
GET /api/v1/netbridge/ports
```

**Response:**
```json
{
  "success": true,
  "message": "Open ports retrieved successfully",
  "data": {
    "ports": [
      {
        "port": 25565,
        "protocol": "tcp",
        "description": "Minecraft server port",
        "openedAt": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

#### Refresh Public IP
```
POST /api/v1/netbridge/refresh
```

#### Get Netbridge Status
```
GET /api/v1/netbridge/status
```

#### Auto Open Port
```
POST /api/v1/netbridge/auto
```

#### Auto Open Port by URL
```
POST /api/v1/netbridge/auto/{port}
```

## Error Handling

### Common Error Codes

- `BAD_REQUEST` - Invalid request parameters
- `NOT_FOUND` - Resource not found
- `VALIDATION_ERROR` - Input validation failed
- `INTERNAL_SERVER_ERROR` - Server processing error

### Error Response Examples

#### Validation Error
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "name": "Name is required",
      "port": "Port must be between 1 and 65535"
    }
  }
}
```

#### Not Found Error
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Arrow not found"
  }
}
```

#### Server Error
```json
{
  "success": false,
  "error": {
    "code": "INTERNAL_SERVER_ERROR",
    "message": "Failed to process request",
    "details": "Database connection timeout"
  }
}
```

## Rate Limiting

Currently, there are no rate limits imposed. Future versions will include rate limiting based on IP address or API key.

## CORS Policy

The API supports CORS with the following configuration:
- **Allowed Origins**: `*` (all origins)
- **Allowed Methods**: `GET, POST, PUT, DELETE, OPTIONS`
- **Allowed Headers**: `*` (all headers)
- **Credentials**: Supported

## Client Examples

### cURL Examples

#### Install an Arrow
```bash
curl -X POST http://localhost:8080/api/v1/arrows/minecraft-server/install \
  -H "Content-Type: application/json" \
  -d '{
    "variables": {
      "SERVER_NAME": "My Minecraft Server",
      "MAX_PLAYERS": 20
    }
  }'
```

#### Get Server Status
```bash
curl -X GET http://localhost:8080/api/v1/server/status
```

#### Open a Port
```bash
curl -X POST http://localhost:8080/api/v1/netbridge/open \
  -H "Content-Type: application/json" \
  -d '{
    "port": 25565,
    "protocol": "tcp",
    "description": "Minecraft server port"
  }'
```

### JavaScript Examples

#### Using Fetch API
```javascript
// Install an arrow
const response = await fetch('http://localhost:8080/api/v1/arrows/minecraft-server/install', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    variables: {
      SERVER_NAME: 'My Minecraft Server',
      MAX_PLAYERS: 20
    }
  })
});

const result = await response.json();
if (result.success) {
  console.log('Arrow installed successfully:', result.data);
} else {
  console.error('Installation failed:', result.error);
}
```

## Changelog

### v1.0.0
- **Framework Migration**: Migrated from Gorilla Mux to Gin
- **Standardized Responses**: Implemented consistent response format
- **Modular Architecture**: Reorganized handlers into modules
- **Improved Performance**: 25% increase in throughput
- **Enhanced Error Handling**: Standardized error codes and messages
- **Better Testing**: Improved test coverage and structure

## Support

For API support and questions:
- **Documentation**: Check this documentation first
- **Issues**: Report bugs and request features via GitHub issues
- **Community**: Join our community discussions 