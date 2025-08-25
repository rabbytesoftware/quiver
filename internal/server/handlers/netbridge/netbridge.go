package netbridge

import (
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/netbridge"
	"github.com/rabbytesoftware/quiver/internal/server/response"
)

// Handler handles netbridge-related HTTP requests
type Handler struct {
	netbridge *netbridge.Netbridge
	logger    *logger.Logger
}

// NewHandler creates a new netbridge handler
func NewHandler(nb *netbridge.Netbridge, logger *logger.Logger) *Handler {
	return &Handler{
		netbridge: nb,
		logger:    logger,
	}
}

// OpenPortRequest represents the request body for opening a port
type OpenPortRequest struct {
	Port     uint16 `json:"port" binding:"required,min=1,max=65535"`
	Protocol string `json:"protocol"`
}

// ClosePortRequest represents the request body for closing a port
type ClosePortRequest struct {
	Port     uint16 `json:"port" binding:"required,min=1,max=65535"`
	Protocol string `json:"protocol"`
}

// OpenPortAutoRequest represents the request body for opening a port automatically
type OpenPortAutoRequest struct {
	Protocol string `json:"protocol"`
}

// PortInfo represents information about an open port
type PortInfo struct {
	Name     string `json:"name"`
	Port     uint16 `json:"port"`
	Host     string `json:"host"`
	Protocol string `json:"protocol"`
}

// OpenPort handles POST /api/v1/netbridge/open
func (h *Handler) OpenPort(c *gin.Context) {
	h.logger.Debug("Handling netbridge open port request")

	var req OpenPortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to decode open port request: %v", err)
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Set default protocol if not provided
	if req.Protocol == "" {
		req.Protocol = "tcp"
	}

	h.logger.Info("Opening port %d with protocol %s", req.Port, req.Protocol)

	result, err := h.netbridge.OpenPort(req.Port, req.Protocol)
	if err != nil {
		h.logger.Error("Failed to open port: %v", err)
		response.BadRequest(c, "Failed to open port", err.Error())
		return
	}

	h.logger.Info("Port %d opened successfully using %s", req.Port, result.Method)
	
	responseData := gin.H{
		"port":     req.Port,
		"protocol": req.Protocol,
		"result":   result,
	}

	response.Success(c, "Port opened successfully", responseData)
}

// ClosePort handles POST /api/v1/netbridge/close
func (h *Handler) ClosePort(c *gin.Context) {
	h.logger.Debug("Handling netbridge close port request")

	var req ClosePortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to decode close port request: %v", err)
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Set default protocol if not provided
	if req.Protocol == "" {
		req.Protocol = "tcp"
	}

	h.logger.Info("Closing port %d with protocol %s", req.Port, req.Protocol)

	result, err := h.netbridge.ClosePort(req.Port, req.Protocol)
	if err != nil {
		h.logger.Error("Failed to close port: %v", err)
		response.BadRequest(c, "Failed to close port", err.Error())
		return
	}

	h.logger.Info("Port %d closed successfully using %s", req.Port, result.Method)
	
	responseData := gin.H{
		"port":     req.Port,
		"protocol": req.Protocol,
		"result":   result,
	}

	response.Success(c, "Port closed successfully", responseData)
}

// OpenPortByURL handles POST /api/v1/netbridge/open/{port}/{protocol}
func (h *Handler) OpenPortByURL(c *gin.Context) {
	h.logger.Debug("Handling netbridge open port by URL request")

	portStr := c.Param("port")
	protocol := c.Param("protocol")
	
	// URL decode the protocol in case it contains encoded characters
	if decodedProtocol, err := url.QueryUnescape(protocol); err == nil {
		protocol = decodedProtocol
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil || port == 0 || port > 65535 {
		response.BadRequest(c, "Invalid port number")
		return
	}

	if protocol == "" {
		protocol = "tcp"
	}

	h.logger.Info("Opening port %d with protocol %s", port, protocol)

	result, err := h.netbridge.OpenPort(uint16(port), protocol)
	if err != nil {
		h.logger.Error("Failed to open port: %v", err)
		response.BadRequest(c, "Failed to open port", err.Error())
		return
	}

	h.logger.Info("Port %d opened successfully using %s", port, result.Method)
	
	responseData := gin.H{
		"port":     port,
		"protocol": protocol,
		"result":   result,
	}

	response.Success(c, "Port opened successfully", responseData)
}

// ClosePortByURL handles POST /api/v1/netbridge/close/{port}/{protocol}
func (h *Handler) ClosePortByURL(c *gin.Context) {
	h.logger.Debug("Handling netbridge close port by URL request")

	portStr := c.Param("port")
	protocol := c.Param("protocol")
	
	// URL decode the protocol in case it contains encoded characters
	if decodedProtocol, err := url.QueryUnescape(protocol); err == nil {
		protocol = decodedProtocol
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil || port == 0 || port > 65535 {
		response.BadRequest(c, "Invalid port number")
		return
	}

	if protocol == "" {
		protocol = "tcp"
	}

	h.logger.Info("Closing port %d with protocol %s", port, protocol)

	result, err := h.netbridge.ClosePort(uint16(port), protocol)
	if err != nil {
		h.logger.Error("Failed to close port: %v", err)
		response.BadRequest(c, "Failed to close port", err.Error())
		return
	}

	h.logger.Info("Port %d closed successfully using %s", port, result.Method)
	
	responseData := gin.H{
		"port":     port,
		"protocol": protocol,
		"result":   result,
	}

	response.Success(c, "Port closed successfully", responseData)
}

// ListOpenPorts handles GET /api/v1/netbridge/ports
func (h *Handler) ListOpenPorts(c *gin.Context) {
	h.logger.Debug("Handling list open ports request")

	if h.netbridge == nil {
		response.InternalServerError(c, "Netbridge not available")
		return
	}

	openPorts := h.netbridge.ListOpenPorts()
	publicIP := h.netbridge.GetPublicIP()

	// Convert to PortInfo format
	var ports []PortInfo
	for _, port := range openPorts {
		ports = append(ports, PortInfo{
			Name:     port.Name(),
			Port:     port.Port(),
			Host:     port.Host(),
			Protocol: port.Protocol(),
		})
	}

	responseData := gin.H{
		"ports":     ports,
		"count":     len(ports),
		"public_ip": publicIP.String(),
	}

	response.Success(c, "Open ports retrieved successfully", responseData)
}

// RefreshPublicIP handles POST /api/v1/netbridge/refresh-ip
func (h *Handler) RefreshPublicIP(c *gin.Context) {
	h.logger.Debug("Handling refresh public IP request")

	if h.netbridge == nil {
		response.InternalServerError(c, "Netbridge not available")
		return
	}

	err := h.netbridge.RefreshPublicIP()
	if err != nil {
		h.logger.Error("Failed to refresh public IP: %v", err)
		response.InternalServerError(c, "Failed to refresh public IP", err.Error())
		return
	}

	publicIP := h.netbridge.GetPublicIP()
	responseData := gin.H{
		"public_ip": publicIP.String(),
	}

	response.Success(c, "Public IP refreshed successfully", responseData)
}

// GetStatus handles GET /api/v1/netbridge/status
func (h *Handler) GetStatus(c *gin.Context) {
	h.logger.Debug("Handling netbridge status request")

	if h.netbridge == nil {
		responseData := gin.H{
			"status":    "disabled",
			"available": false,
			"reason":    "Netbridge not initialized",
		}
		response.Success(c, "Netbridge status retrieved", responseData)
		return
	}

	openPorts := h.netbridge.ListOpenPorts()
	publicIP := h.netbridge.GetPublicIP()

	responseData := gin.H{
		"status":        "enabled",
		"available":     true,
		"open_ports":    len(openPorts),
		"public_ip":     publicIP.String(),
		"methods":       []string{"upnp", "natpmp", "manual"},
	}

	response.Success(c, "Netbridge status retrieved successfully", responseData)
}

// OpenPortAuto handles POST /api/v1/netbridge/open-auto
func (h *Handler) OpenPortAuto(c *gin.Context) {
	h.logger.Debug("Handling netbridge open port auto request")

	var req OpenPortAutoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Protocol is optional, so we'll set a default
		req.Protocol = "tcp"
	}

	if req.Protocol == "" {
		req.Protocol = "tcp"
	}

	h.logger.Info("Opening port automatically with protocol %s", req.Protocol)

	result, err := h.netbridge.OpenPortAuto(req.Protocol)
	if err != nil {
		h.logger.Error("Failed to open port automatically: %v", err)
		response.BadRequest(c, "Failed to open port automatically", err.Error())
		return
	}

	h.logger.Info("Port %d opened automatically using %s", result.Port, result.Method)
	
	responseData := gin.H{
		"protocol": req.Protocol,
		"result":   result,
	}

	response.Success(c, "Port opened automatically", responseData)
}

// OpenPortAutoByURL handles POST /api/v1/netbridge/open-auto/{protocol}
func (h *Handler) OpenPortAutoByURL(c *gin.Context) {
	h.logger.Debug("Handling netbridge open port auto by URL request")

	protocol := c.Param("protocol")
	if protocol == "" {
		protocol = "tcp"
	}

	// URL decode the protocol in case it contains encoded characters
	if decodedProtocol, err := url.QueryUnescape(protocol); err == nil {
		protocol = decodedProtocol
	}

	h.logger.Info("Opening port automatically with protocol %s", protocol)

	result, err := h.netbridge.OpenPortAuto(protocol)
	if err != nil {
		h.logger.Error("Failed to open port automatically: %v", err)
		response.BadRequest(c, "Failed to open port automatically", err.Error())
		return
	}

	h.logger.Info("Port %d opened automatically using %s", result.Port, result.Method)
	
	responseData := gin.H{
		"protocol": protocol,
		"result":   result,
	}

	response.Success(c, "Port opened automatically", responseData)
} 