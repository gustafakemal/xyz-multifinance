package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/username/xyz-multifinance/internal/domain"
	"github.com/username/xyz-multifinance/internal/infrastructure/security"
	"github.com/username/xyz-multifinance/internal/usecase"
)

// TransactionHandler handles HTTP requests for transactions
type TransactionHandler struct {
	usecase *usecase.TransactionUsecase
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(usecase *usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{usecase: usecase}
}

// CreateTransaction handles transaction creation requests
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		security.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request body
	var req usecase.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		security.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Input validation to prevent injection attacks (OWASP #3)
	if req.ConsumerID <= 0 {
		security.WriteError(w, http.StatusBadRequest, "Invalid consumer ID")
		return
	}

	if len(req.AssetName) == 0 || len(req.AssetName) > 100 {
		security.WriteError(w, http.StatusBadRequest, "Asset name must be between 1 and 100 characters")
		return
	}

	// Create transaction
	transaction, err := h.usecase.CreateTransaction(req)
	if err != nil {
		switch err {
		case domain.ErrConsumerNotFound:
			security.WriteError(w, http.StatusNotFound, "Consumer not found")
		case domain.ErrLimitNotFound:
			security.WriteError(w, http.StatusNotFound, "Consumer limit not found for specified tenor")
		case domain.ErrLimitExceeded:
			security.WriteError(w, http.StatusBadRequest, "Insufficient credit limit")
		case domain.ErrInvalidTenor, domain.ErrInvalidOTR, domain.ErrInvalidAdminFee,
			domain.ErrInvalidInstallment, domain.ErrInvalidAssetName:
			security.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			security.WriteError(w, http.StatusInternalServerError, "Failed to create transaction")
		}
		return
	}

	// Return success response
	security.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":     "Transaction created successfully",
		"transaction": transaction,
	})
}

// GetTransaction retrieves a transaction by contract number
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		security.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get contract number from query parameter
	contractNumber := r.URL.Query().Get("contract_number")
	if contractNumber == "" {
		security.WriteError(w, http.StatusBadRequest, "Contract number is required")
		return
	}

	// Input sanitization (OWASP #3)
	if len(contractNumber) > 100 {
		security.WriteError(w, http.StatusBadRequest, "Invalid contract number")
		return
	}

	transaction, err := h.usecase.GetTransactionByContract(contractNumber)
	if err != nil {
		if err == domain.ErrTransactionNotFound {
			security.WriteError(w, http.StatusNotFound, "Transaction not found")
		} else {
			security.WriteError(w, http.StatusInternalServerError, "Failed to retrieve transaction")
		}
		return
	}

	security.WriteJSON(w, http.StatusOK, transaction)
}

// GetConsumerTransactions retrieves all transactions for a consumer
func (h *TransactionHandler) GetConsumerTransactions(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		security.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get consumer ID from query parameter
	consumerIDStr := r.URL.Query().Get("consumer_id")
	if consumerIDStr == "" {
		security.WriteError(w, http.StatusBadRequest, "Consumer ID is required")
		return
	}

	consumerID, err := strconv.ParseInt(consumerIDStr, 10, 64)
	if err != nil || consumerID <= 0 {
		security.WriteError(w, http.StatusBadRequest, "Invalid consumer ID")
		return
	}

	transactions, err := h.usecase.GetConsumerTransactions(consumerID)
	if err != nil {
		if err == domain.ErrConsumerNotFound {
			security.WriteError(w, http.StatusNotFound, "Consumer not found")
		} else {
			security.WriteError(w, http.StatusInternalServerError, "Failed to retrieve transactions")
		}
		return
	}

	security.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"consumer_id":  consumerID,
		"transactions": transactions,
	})
}

// GetConsumerLimits retrieves all limits for a consumer
func (h *TransactionHandler) GetConsumerLimits(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		security.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get consumer ID from query parameter
	consumerIDStr := r.URL.Query().Get("consumer_id")
	if consumerIDStr == "" {
		security.WriteError(w, http.StatusBadRequest, "Consumer ID is required")
		return
	}

	consumerID, err := strconv.ParseInt(consumerIDStr, 10, 64)
	if err != nil || consumerID <= 0 {
		security.WriteError(w, http.StatusBadRequest, "Invalid consumer ID")
		return
	}

	limits, err := h.usecase.GetConsumerLimits(consumerID)
	if err != nil {
		if err == domain.ErrConsumerNotFound {
			security.WriteError(w, http.StatusNotFound, "Consumer not found")
		} else {
			security.WriteError(w, http.StatusInternalServerError, "Failed to retrieve limits")
		}
		return
	}

	security.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"consumer_id": consumerID,
		"limits":      limits,
	})
}

// HealthCheck provides a health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	security.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   "2026-02-01",
	})
}
