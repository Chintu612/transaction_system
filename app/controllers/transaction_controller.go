package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"transaction_system/app/models"
	"transaction_system/app/repositories"
	"transaction_system/app/services"

	"github.com/julienschmidt/httprouter"
)

type TransactionControllerI interface {
	CreateTransaction(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	GetTransactionsByType(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	GetTransitiveSum(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

type transactionController struct {
	transactionService services.TransactionServiceI
}

func NewTransactionController() TransactionControllerI {
	return &transactionController{
		transactionService: services.NewTransactionService(),
	}
}

func MakeTransactionController(transactionService services.TransactionServiceI) TransactionControllerI {
	return &transactionController{
		transactionService: transactionService,
	}
}

func (t *transactionController) CreateTransaction(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Retrieve transaction ID from URL params
	transactionID := params.ByName("transaction_id")
	if transactionID == "" {
		respondWithError(w, "Transaction ID is required", http.StatusBadRequest)
		return
	}

	// Convert transaction ID to uint
	transactionIDUint, err := strconv.ParseUint(transactionID, 10, 64)
	if err != nil {
		respondWithError(w, "Invalid transaction ID format", http.StatusBadRequest)
		return
	}

	// Decode request body
	var transactionData map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&transactionData); err != nil {
		respondWithError(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Validate schema
	validatedData, isValid := validateSchema(w, transactionData, "amount", "type")

	if !isValid {
		return
	}

	// Extract transaction details from parsed data
	amount, ok := validatedData["amount"].(float64)
	if !ok {
		respondWithError(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	transactionType, ok := validatedData["type"].(string)
	if !ok {
		respondWithError(w, "Invalid type format", http.StatusBadRequest)
		return
	}

	// Extract parent_id and set it to nil if not present
	var parentID *uint
	if val, exists := transactionData["parent_id"]; exists {
		if parentIDValue, isUint := val.(float64); isUint {
			parentIDValueUint := uint(parentIDValue)
			parentID = &parentIDValueUint
		} else {
			respondWithError(w, "Invalid parent_id format", http.StatusBadRequest)
			return
		}
	}

	// Create a new transaction object
	newTransaction := models.Transaction{
		Id:       uint(transactionIDUint),
		Amount:   amount,
		Type:     transactionType,
		ParentID: parentID,
	}

	// Call the service to create the transaction
	status, err := t.transactionService.CreateTransaction(newTransaction)
	if err != nil {
		if err == services.ErrParentTransactionNotFound {
			// Handling "Parent transaction does not exist" as Bad Request
			respondWithError(w, "Parent transaction does not exist", http.StatusBadRequest)
			return
		}
		if err == repositories.ErrTransactionAlreadyExist {
			// Handling "transaction does not exist" as Bad Request
			respondWithError(w, "transaction with the same ID already exists", http.StatusBadRequest)
			return
		}
		respondWithError(w, fmt.Sprintf("Error creating transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the created transaction status
	response := map[string]string{
		"status": getStatusMessage(status),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		respondWithError(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(jsonResponse)
	if err != nil {
		respondWithError(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

func (t *transactionController) GetTransactionsByType(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Extract transaction type from URL params
	transactionType := params.ByName("type")
	if transactionType == "" {
		http.Error(w, "Transaction type is required", http.StatusBadRequest)
		return
	}

	// Call the service to get transaction IDs by type
	transactionIDs, err := t.transactionService.GetTransactionIDsByType(transactionType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving transaction IDs: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the list of transaction IDs
	response := map[string][]uint{"transaction_ids": transactionIDs}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		respondWithError(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		respondWithError(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

func (t *transactionController) GetTransitiveSum(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	transactionID := params.ByName("transaction_id")
	if transactionID == "" {
		respondWithError(w, "Transaction ID is required", http.StatusBadRequest)
		return
	}

	// Convert transactionID to uint
	transactionIDUint, err := strconv.ParseUint(transactionID, 10, 64)
	if err != nil {
		respondWithError(w, "Invalid transaction ID format", http.StatusBadRequest)
		return
	}

	// Call the service to get the sum
	sum, err := t.transactionService.GetTransitiveSum(uint(transactionIDUint))
	if err != nil {
		if err == services.ErrTransactionNotFound {
			// Handling "Transaction does not exist" as Bad Request
			respondWithError(w, "Transaction does not exist for given transaction ID", http.StatusBadRequest)
			return
		}
		respondWithError(w, fmt.Sprintf("Error retrieving transitive sum: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the sum
	response := map[string]float64{"sum": sum}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		respondWithError(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		respondWithError(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

// getStatusMessage returns a human-readable status message based on the transaction creation status.
func getStatusMessage(status bool) string {
	if status {
		return "ok"
	}
	return "unable to create transaction"
}

func respondWithError(w http.ResponseWriter, errMsg string, statusCode int) {
	errorResponse := map[string]interface{}{
		"success": "false",
		"error":   errMsg,
		"status":  statusCode,
	}

	jsonResponse, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

func validateSchema(w http.ResponseWriter, data map[string]interface{}, expectedFields ...string) (map[string]interface{}, bool) {
	result := make(map[string]interface{})

	for _, field := range expectedFields {
		val, exists := data[field]
		if !exists {
			respondWithError(w, fmt.Sprintf("Field '%s' is missing", field), http.StatusBadRequest)
			return nil, false
		}
		result[field] = val
	}

	return result, true
}
