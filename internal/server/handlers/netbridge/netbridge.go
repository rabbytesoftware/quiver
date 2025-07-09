package netbridge

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
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
	Port     uint16 `json:"port"`
	Protocol string `json:"protocol"`
}

// ClosePortRequest represents the request body for closing a port
type ClosePortRequest struct {
	Port     uint16 `json:"port"`
	Protocol string `json:"protocol"`
}

// PortStatusResponse represents the response for port operations
type PortStatusResponse struct {
	Success bool                                `json:"success"`
	Result  *netbridge.PortForwardingResult     `json:"result,omitempty"`
	Error   string                              `json:"error,omitempty"`
}

// OpenPortsListResponse represents the response for listing open ports
type OpenPortsListResponse struct {
	Success   bool                  `json:"success"`
	Ports     []PortInfo            `json:"ports"`
	PublicIP  string                `json:"public_ip"`
	Error     string                `json:"error,omitempty"`
}

// PortInfo represents information about an open port
type PortInfo struct {
	Name     string `json:"name"`
	Port     uint16 `json:"port"`
	Host     string `json:"host"`
	Protocol string `json:"protocol"`
}

// OpenPortHandler handles POST /api/v1/netbridge/open
func (h *Handler) OpenPortHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling netbridge open port request")

	var req OpenPortRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode open port request: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate port number
	if req.Port == 0 || req.Port > 65535 {
		response.WriteError(w, http.StatusBadRequest, "Port must be between 1 and 65535")
		return
	}

	// Validate protocol
	if req.Protocol == "" {
		req.Protocol = "tcp" // Default to TCP
	}

	h.logger.Info("Opening port %d with protocol %s", req.Port, req.Protocol)

	result, err := h.netbridge.OpenPort(req.Port, req.Protocol)
	if err != nil {
		h.logger.Error("Failed to open port: %v", err)
		response.WriteJSON(w, http.StatusBadRequest, PortStatusResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info("Port %d opened successfully using %s", req.Port, result.Method)
	response.WriteJSON(w, http.StatusOK, PortStatusResponse{
		Success: true,
		Result:  result,
	})
}

// ClosePortHandler handles POST /api/v1/netbridge/close
func (h *Handler) ClosePortHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling netbridge close port request")

	var req ClosePortRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode close port request: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate port number
	if req.Port == 0 || req.Port > 65535 {
		response.WriteError(w, http.StatusBadRequest, "Port must be between 1 and 65535")
		return
	}

	// Validate protocol
	if req.Protocol == "" {
		req.Protocol = "tcp" // Default to TCP
	}

	h.logger.Info("Closing port %d with protocol %s", req.Port, req.Protocol)

	result, err := h.netbridge.ClosePort(req.Port, req.Protocol)
	if err != nil {
		h.logger.Error("Failed to close port: %v", err)
		response.WriteJSON(w, http.StatusBadRequest, PortStatusResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info("Port %d closed successfully using %s", req.Port, result.Method)
	response.WriteJSON(w, http.StatusOK, PortStatusResponse{
		Success: true,
		Result:  result,
	})
}

// OpenPortByURLHandler handles POST /api/v1/netbridge/open/{port}/{protocol}
func (h *Handler) OpenPortByURLHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling netbridge open port by URL request")

	vars := mux.Vars(r)
	portStr := vars["port"]
	protocol := vars["protocol"]
	
	// URL decode the protocol in case it contains encoded characters like %2F for /
	if decodedProtocol, err := url.QueryUnescape(protocol); err == nil {
		protocol = decodedProtocol
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil || port == 0 || port > 65535 {
		response.WriteError(w, http.StatusBadRequest, "Invalid port number")
		return
	}

	if protocol == "" {
		protocol = "tcp"
	}

	h.logger.Info("Opening port %d with protocol %s", port, protocol)

	result, err := h.netbridge.OpenPort(uint16(port), protocol)
	if err != nil {
		h.logger.Error("Failed to open port: %v", err)
		response.WriteJSON(w, http.StatusBadRequest, PortStatusResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info("Port %d opened successfully using %s", port, result.Method)
	response.WriteJSON(w, http.StatusOK, PortStatusResponse{
		Success: true,
		Result:  result,
	})
}

// ClosePortByURLHandler handles POST /api/v1/netbridge/close/{port}/{protocol}
func (h *Handler) ClosePortByURLHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling netbridge close port by URL request")

	vars := mux.Vars(r)
	portStr := vars["port"]
	protocol := vars["protocol"]
	
	// URL decode the protocol in case it contains encoded characters like %2F for /
	if decodedProtocol, err := url.QueryUnescape(protocol); err == nil {
		protocol = decodedProtocol
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil || port == 0 || port > 65535 {
		response.WriteError(w, http.StatusBadRequest, "Invalid port number")
		return
	}

	if protocol == "" {
		protocol = "tcp"
	}

	h.logger.Info("Closing port %d with protocol %s", port, protocol)

	result, err := h.netbridge.ClosePort(uint16(port), protocol)
	if err != nil {
		h.logger.Error("Failed to close port: %v", err)
		response.WriteJSON(w, http.StatusBadRequest, PortStatusResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info("Port %d closed successfully using %s", port, result.Method)
	response.WriteJSON(w, http.StatusOK, PortStatusResponse{
		Success: true,
		Result:  result,
	})
}

// ListOpenPortsHandler handles GET /api/v1/netbridge/ports
func (h *Handler) ListOpenPortsHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling netbridge list open ports request")

	ports := h.netbridge.ListOpenPorts()
	publicIP := h.netbridge.GetPublicIP()

	// Convert ports to response format
	portInfos := make([]PortInfo, len(ports))
	for i, port := range ports {
		portInfos[i] = PortInfo{
			Name:     port.Name(),
			Port:     port.Port(),
			Host:     port.Host(),
			Protocol: port.Protocol(),
		}
	}

	h.logger.Info("Retrieved %d open ports", len(ports))
	response.WriteJSON(w, http.StatusOK, OpenPortsListResponse{
		Success:  true,
		Ports:    portInfos,
		PublicIP: publicIP.String(),
	})
}

// RefreshPublicIPHandler handles POST /api/v1/netbridge/refresh-ip
func (h *Handler) RefreshPublicIPHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling netbridge refresh public IP request")

	if err := h.netbridge.RefreshPublicIP(); err != nil {
		h.logger.Error("Failed to refresh public IP: %v", err)
		response.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	publicIP := h.netbridge.GetPublicIP()
	h.logger.Info("Public IP refreshed: %s", publicIP.String())
	
	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"public_ip": publicIP.String(),
	})
}

// GetStatusHandler handles GET /api/v1/netbridge/status
func (h *Handler) GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling netbridge status request")

	ports := h.netbridge.ListOpenPorts()
	publicIP := h.netbridge.GetPublicIP()

	// Convert ports to response format
	portInfos := make([]PortInfo, len(ports))
	for i, port := range ports {
		portInfos[i] = PortInfo{
			Name:     port.Name(),
			Port:     port.Port(),
			Host:     port.Host(),
			Protocol: port.Protocol(),
		}
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success":       true,
		"public_ip":     publicIP.String(),
		"open_ports":    portInfos,
		"total_ports":   len(ports),
		"methods":       []string{"upnp", "manual", "natpmp"},
		"available":     true,
	})
} 

// OpenPortAutoRequest represents the request body for opening an automatic port
type OpenPortAutoRequest struct {
	Protocol string `json:"protocol"`
}

// OpenPortAutoHandler handles POST /api/v1/netbridge/open-auto
func (h *Handler) OpenPortAutoHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling netbridge open auto port request")

	var req OpenPortAutoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode open auto port request: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate protocol
	if req.Protocol == "" {
		req.Protocol = "tcp" // Default to TCP
	}

	h.logger.Info("Opening automatic port with protocol %s", req.Protocol)

	result, err := h.netbridge.OpenPortAuto(req.Protocol)
	if err != nil {
		h.logger.Error("Failed to open automatic port: %v", err)
		response.WriteJSON(w, http.StatusBadRequest, PortStatusResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info("Automatic port opened successfully: %d using %s", result.Port.Port(), result.Method)
	response.WriteJSON(w, http.StatusOK, PortStatusResponse{
		Success: true,
		Result:  result,
	})
}

// OpenPortAutoByURLHandler handles POST /api/v1/netbridge/open-auto/{protocol}
func (h *Handler) OpenPortAutoByURLHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling netbridge open auto port by URL request")

	vars := mux.Vars(r)
	protocol := vars["protocol"]
	
	// URL decode the protocol in case it contains encoded characters like %2F for /
	if decodedProtocol, err := url.QueryUnescape(protocol); err == nil {
		protocol = decodedProtocol
	}

	if protocol == "" {
		protocol = "tcp"
	}

	h.logger.Info("Opening automatic port with protocol %s", protocol)

	result, err := h.netbridge.OpenPortAuto(protocol)
	if err != nil {
		h.logger.Error("Failed to open automatic port: %v", err)
		response.WriteJSON(w, http.StatusBadRequest, PortStatusResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info("Automatic port opened successfully: %d using %s", result.Port.Port(), result.Method)
	response.WriteJSON(w, http.StatusOK, PortStatusResponse{
		Success: true,
		Result:  result,
	})
} 