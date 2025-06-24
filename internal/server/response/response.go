package response

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes a JSON response with the given status code and data
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// WriteError writes an error response with the given status code and message
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	errorResponse := map[string]string{
		"error": message,
	}
	WriteJSON(w, statusCode, errorResponse)
}

// WriteSuccess writes a success response with the given message
func WriteSuccess(w http.ResponseWriter, message string) {
	successResponse := map[string]string{
		"message": message,
		"status":  "success",
	}
	WriteJSON(w, http.StatusOK, successResponse)
}

// WriteNotFound writes a 404 not found response
func WriteNotFound(w http.ResponseWriter, resource string) {
	WriteError(w, http.StatusNotFound, resource+" not found")
}

// WriteBadRequest writes a 400 bad request response
func WriteBadRequest(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusBadRequest, message)
}

// WriteInternalServerError writes a 500 internal server error response
func WriteInternalServerError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusInternalServerError, message)
}

// WriteUnauthorized writes a 401 unauthorized response
func WriteUnauthorized(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusUnauthorized, message)
}

// WriteForbidden writes a 403 forbidden response
func WriteForbidden(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusForbidden, message)
} 