package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"transaction_system/app/controllers"
	"transaction_system/app/repositories"
	"transaction_system/app/services"
	"transaction_system/app/services/mock_services"
)

func TestCreateTransaction_Success(t *testing.T) {
	t.Run("when parent_id is not provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Mocks
		mockTransactionService := mock_services.NewMockTransactionServiceI(ctrl)

		// Controller
		transactionController := controllers.MakeTransactionController(mockTransactionService)

		//transactionID
		transactionID := 1

		// Request body
		requestBody := map[string]interface{}{
			"amount": 100.0,
			"type":   "purchase",
		}

		// Convert request body to JSON
		jsonRequest, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatal(err)
		}

		// Mock expectations
		mockTransactionService.EXPECT().CreateTransaction(gomock.Any()).Return(true, nil)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/transactionservice/transaction/%d", transactionID), bytes.NewBuffer(jsonRequest))
		recorder := httptest.NewRecorder()
		router := httprouter.New()
		router.Handle(http.MethodPut, "/transactionservice/transaction/:transaction_id", transactionController.CreateTransaction)
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusCreated, recorder.Code)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, "ok", resData["status"])
		assert.NotContains(t, resData, "error") // No error key should be present
	})

	t.Run("when parent_id is provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Mocks
		mockTransactionService := mock_services.NewMockTransactionServiceI(ctrl)

		// Controller
		transactionController := controllers.MakeTransactionController(mockTransactionService)

		//transactionID
		transactionID := 3

		// Request body
		requestBody := map[string]interface{}{
			"amount":    100.0,
			"type":      "purchase",
			"parent_id": 1,
		}

		// Convert request body to JSON
		jsonRequest, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatal(err)
		}

		// Mock expectations
		mockTransactionService.EXPECT().CreateTransaction(gomock.Any()).Return(true, nil)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/transactionservice/transaction/%d", transactionID), bytes.NewBuffer(jsonRequest))
		recorder := httptest.NewRecorder()
		router := httprouter.New()
		router.Handle(http.MethodPut, "/transactionservice/transaction/:transaction_id", transactionController.CreateTransaction)
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusCreated, recorder.Code)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, "ok", resData["status"])
		assert.NotContains(t, resData, "error") // No error key should be present
	})
}

func TestCreateTransaction_InvalidTransactionID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks
	mockTransactionService := mock_services.NewMockTransactionServiceI(ctrl)

	// Controller
	transactionController := controllers.MakeTransactionController(mockTransactionService)

	// Request body
	requestBody := map[string]interface{}{
		"amount": 100.0,
		"type":   "purchase",
	}

	// Convert request body to JSON
	jsonRequest, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/transactionservice/transaction/invalid_id", bytes.NewBuffer(jsonRequest))
	recorder := httptest.NewRecorder()
	router := httprouter.New()
	router.Handle(http.MethodPut, "/transactionservice/transaction/:transaction_id", transactionController.CreateTransaction)
	router.ServeHTTP(recorder, req)

	// Assert status code is BadRequest
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Assert response body contains expected error message
	expectedResponse := `{"error":"Invalid transaction ID format","status":400,"success":"false"}`
	assert.Equal(t, expectedResponse, recorder.Body.String())
}

func TestCreateTransaction_InvalidAmountFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks
	mockTransactionService := mock_services.NewMockTransactionServiceI(ctrl)

	// Controller
	transactionController := controllers.MakeTransactionController(mockTransactionService)

	//transactionID
	transactionID := 1

	// Request body with invalid amount format (string instead of float)
	requestBody := map[string]interface{}{
		"amount": "invalid_amount",
		"type":   "purchase",
	}

	// Convert request body to JSON
	jsonRequest, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/transactionservice/transaction/%d", transactionID), bytes.NewBuffer(jsonRequest))
	recorder := httptest.NewRecorder()
	router := httprouter.New()
	router.Handle(http.MethodPut, "/transactionservice/transaction/:transaction_id", transactionController.CreateTransaction)
	router.ServeHTTP(recorder, req)

	// Assert status code is BadRequest
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Assert response body contains expected error message
	expectedResponse := `{"error":"Invalid amount format","status":400,"success":"false"}`
	assert.Equal(t, expectedResponse, recorder.Body.String())
}

func TestCreateTransaction_MissingAmountField(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks
	mockTransactionService := mock_services.NewMockTransactionServiceI(ctrl)

	// Controller
	transactionController := controllers.MakeTransactionController(mockTransactionService)

	//transactionID
	transactionID := 1

	// Request body without the amount field
	requestBody := map[string]interface{}{
		"type": "purchase",
	}

	// Convert request body to JSON
	jsonRequest, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/transactionservice/transaction/%d", transactionID), bytes.NewBuffer(jsonRequest))
	recorder := httptest.NewRecorder()
	router := httprouter.New()
	router.Handle(http.MethodPut, "/transactionservice/transaction/:transaction_id", transactionController.CreateTransaction)
	router.ServeHTTP(recorder, req)

	// Assert status code is BadRequest
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Assert response body contains expected error message
	expectedResponse := `{"error":"Field 'amount' is missing","status":400,"success":"false"}`
	assert.Equal(t, expectedResponse, recorder.Body.String())
}

func TestCreateTransaction_InvalidParentIDFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks
	mockTransactionService := mock_services.NewMockTransactionServiceI(ctrl)

	// Controller
	transactionController := controllers.MakeTransactionController(mockTransactionService)

	// Transaction ID
	transactionID := "123"

	// Request body with invalid parent_id format (non-numeric)
	requestBody := map[string]interface{}{
		"amount":    100.0,
		"type":      "purchase",
		"parent_id": "invalid_parent_id",
	}

	// Convert request body to JSON
	jsonRequest, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/transactionservice/transaction/%s", transactionID), bytes.NewBuffer(jsonRequest))
	recorder := httptest.NewRecorder()
	router := httprouter.New()
	router.Handle(http.MethodPut, "/transactionservice/transaction/:transaction_id", transactionController.CreateTransaction)
	router.ServeHTTP(recorder, req)

	// Assert status code is BadRequest
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Assert response body contains expected error message
	expectedResponse := `{"error":"Invalid parent_id format","status":400,"success":"false"}`
	assert.Equal(t, expectedResponse, recorder.Body.String())
}

func TestCreateTransaction_ParentTransactionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks
	mockTransactionService := mock_services.NewMockTransactionServiceI(ctrl)

	// Controller
	transactionController := controllers.MakeTransactionController(mockTransactionService)

	//transactionID
	transactionID := 1

	// Request body
	requestBody := map[string]interface{}{
		"amount": 100.0,
		"type":   "purchase",
	}

	// Convert request body to JSON
	jsonRequest, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	// Mock expectations
	mockTransactionService.EXPECT().CreateTransaction(gomock.Any()).Return(false, services.ErrParentTransactionNotFound)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/transactionservice/transaction/%d", transactionID), bytes.NewBuffer(jsonRequest))
	recorder := httptest.NewRecorder()
	router := httprouter.New()
	router.Handle(http.MethodPut, "/transactionservice/transaction/:transaction_id", transactionController.CreateTransaction)
	router.ServeHTTP(recorder, req)

	// Assert status code is BadRequest
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Assert response body contains expected error message
	expectedResponse := `{"error":"Parent transaction does not exist","status":400,"success":"false"}`
	assert.Equal(t, expectedResponse, recorder.Body.String())
}

func TestCreateTransaction_TransactionAlreadyExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks
	mockTransactionService := mock_services.NewMockTransactionServiceI(ctrl)

	// Controller
	transactionController := controllers.MakeTransactionController(mockTransactionService)

	//transactionID
	transactionID := 1

	// Request body
	requestBody := map[string]interface{}{
		"amount": 100.0,
		"type":   "purchase",
	}

	// Convert request body to JSON
	jsonRequest, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	// Mock expectations
	mockTransactionService.EXPECT().CreateTransaction(gomock.Any()).Return(false, repositories.ErrTransactionAlreadyExist)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/transactionservice/transaction/%d", transactionID), bytes.NewBuffer(jsonRequest))
	recorder := httptest.NewRecorder()
	router := httprouter.New()
	router.Handle(http.MethodPut, "/transactionservice/transaction/:transaction_id", transactionController.CreateTransaction)
	router.ServeHTTP(recorder, req)

	// Assert status code is BadRequest
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Assert response body contains expected error message
	expectedResponse := `{"error":"transaction with the same ID already exists","status":400,"success":"false"}`
	assert.Equal(t, expectedResponse, recorder.Body.String())
}

func TestCreateTransaction_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks
	mockTransactionService := mock_services.NewMockTransactionServiceI(ctrl)

	// Controller
	transactionController := controllers.MakeTransactionController(mockTransactionService)

	//transactionID
	transactionID := 1

	// Request body
	requestBody := map[string]interface{}{
		"amount": 100.0,
		"type":   "purchase",
	}

	// Convert request body to JSON
	jsonRequest, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	// Mock expectations
	mockTransactionService.EXPECT().CreateTransaction(gomock.Any()).Return(false, errors.New("some internal error"))

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/transactionservice/transaction/%d", transactionID), bytes.NewBuffer(jsonRequest))
	recorder := httptest.NewRecorder()
	router := httprouter.New()
	router.Handle(http.MethodPut, "/transactionservice/transaction/:transaction_id", transactionController.CreateTransaction)
	router.ServeHTTP(recorder, req)

	// Assert status code is InternalServerError
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// Assert response body contains expected error message
	expectedResponse := `{"error":"Error creating transaction: some internal error","status":500,"success":"false"}`
	assert.Equal(t, expectedResponse, recorder.Body.String())
}
