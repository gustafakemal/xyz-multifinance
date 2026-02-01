package handler

import (
	"encoding/json"
	"net/http"

	"github.com/username/xyz-multifinance/internal/infrastructure/security"
	"github.com/username/xyz-multifinance/internal/repository"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	consumerRepo repository.ConsumerRepository
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(consumerRepo repository.ConsumerRepository) *AuthHandler {
	return &AuthHandler{consumerRepo: consumerRepo}
}

// LoginRequest represents login request body
type LoginRequest struct {
	NIK string `json:"nik"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token      string `json:"token"`
	ConsumerID int64  `json:"consumer_id"`
	FullName   string `json:"full_name"`
	ExpiresIn  string `json:"expires_in"`
}

// Login generates JWT token for consumer
// This is a simplified version for testing purposes
// In production, you should add password verification
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		security.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		security.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate NIK
	if len(req.NIK) != 16 {
		security.WriteError(w, http.StatusBadRequest, "NIK must be 16 digits")
		return
	}

	// Find consumer by NIK
	consumer, err := h.consumerRepo.FindByNIK(req.NIK)
	if err != nil {
		security.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := security.GenerateToken(consumer.ID, consumer.NIK)
	if err != nil {
		security.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return token
	response := LoginResponse{
		Token:      token,
		ConsumerID: consumer.ID,
		FullName:   consumer.FullName,
		ExpiresIn:  "24 hours",
	}

	security.WriteJSON(w, http.StatusOK, response)
}
