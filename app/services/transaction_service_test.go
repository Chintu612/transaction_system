package services_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"

	"transaction_system/app/models"
	"transaction_system/app/repositories/mock_repositories"
	"transaction_system/app/services"
)

func TestCreateTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks
	mockTransactionRepo := mock_repositories.NewMockTransactionRepositoryI(ctrl)

	// Service
	transactionService := services.MakeTransactionService(mockTransactionRepo)

	// Test data
	transaction := models.Transaction{
		Id:     1,
		Amount: 100.0,
		Type:   "purchase",
	}

	// Mock expectations
	mockTransactionRepo.EXPECT().Create(gomock.Any()).Return(nil)

	// Test the service method
	status, err := transactionService.CreateTransaction(transaction)

	// Assert the result
	assert.True(t, status)
	assert.NoError(t, err)
}

func TestCreateTransaction_ParentTransactionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks
	mockTransactionRepo := mock_repositories.NewMockTransactionRepositoryI(ctrl)

	// Service
	transactionService := services.MakeTransactionService(mockTransactionRepo)

	parentID := uint(2)

	// Test data
	transaction := models.Transaction{
		Id:       1,
		Amount:   100.0,
		Type:     "purchase",
		ParentID: &parentID,
	}

	// Mock expectations
	mockTransactionRepo.EXPECT().GetByID(parentID).Return(nil, nil) // Set up expectation for GetByID

	// Test the service method
	status, err := transactionService.CreateTransaction(transaction)

	// Assert the result
	assert.False(t, status)
	assert.Error(t, err)
	assert.EqualError(t, err, "parent transaction does not exist")
}

func TestCreateTransaction_TransactionAlreadyExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks
	mockTransactionRepo := mock_repositories.NewMockTransactionRepositoryI(ctrl)

	// Service
	transactionService := services.MakeTransactionService(mockTransactionRepo)

	// Test data
	transaction := models.Transaction{
		Id:     1,
		Amount: 100.0,
		Type:   "purchase",
	}

	// Mock expectations
	mockTransactionRepo.EXPECT().Create(gomock.Any()).Return(errors.New("transaction with the same ID already exists"))

	// Test the service method
	status, err := transactionService.CreateTransaction(transaction)

	// Assert the result
	assert.False(t, status)
	assert.Error(t, err)
	assert.EqualError(t, err, "transaction with the same ID already exists")
}
