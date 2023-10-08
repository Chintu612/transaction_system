package services

import (
	"errors"
	"transaction_system/app/models"
	"transaction_system/app/repositories"
)

//go:generate mockgen -source=./transaction_service.go -destination=mock_services/mock_transaction_service.go -package=mock_services

var ErrParentTransactionNotFound = errors.New("parent transaction does not exist")
var ErrTransactionNotFound = errors.New("transaction does not exist for given transaction ID")

type TransactionServiceI interface {
	CreateTransaction(transaction models.Transaction) (bool, error)
	GetTransactionIDsByType(transactionType string) ([]uint, error)
	GetTransitiveSum(transactionID uint) (float64, error)
}

type transactionService struct {
	transactionRepo repositories.TransactionRepositoryI
}

func NewTransactionService() TransactionServiceI {
	return &transactionService{
		transactionRepo: repositories.NewTransactionRepository(),
	}
}

func MakeTransactionService(transactionRepo repositories.TransactionRepositoryI) TransactionServiceI {
	return &transactionService{
		transactionRepo: transactionRepo,
	}
}

// CreateTransaction creates a new transaction using the provided transaction data.
func (t *transactionService) CreateTransaction(transaction models.Transaction) (bool, error) {

	if transaction.ParentID != nil {
		parentTransaction, err := t.transactionRepo.GetByID(*transaction.ParentID)
		if err != nil {
			return false, err
		}

		if parentTransaction == nil {
			return false, ErrParentTransactionNotFound
		}
	}

	err := t.transactionRepo.Create(&transaction)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetTransactionIDsByType retrieves a list of transaction IDs that match the given transactionType.
func (t *transactionService) GetTransactionIDsByType(transactionType string) ([]uint, error) {
	var transactionIDs []uint

	transactions, err := t.transactionRepo.GetByType(transactionType)
	if err != nil {
		return nil, err
	}

	// Extract transaction IDs from the fetched transactions
	for _, transaction := range transactions {
		transactionIDs = append(transactionIDs, transaction.Id)
	}

	return transactionIDs, nil
}

// GetTransitiveSum retrieves the sum of all transactions transitively linked by their parent_id to a given transaction ID.
func (t *transactionService) GetTransitiveSum(transactionID uint) (float64, error) {
	Transaction, err := t.transactionRepo.GetByID(transactionID)
	if err != nil {
		return 0, err
	}

	if Transaction == nil {
		return 0, ErrTransactionNotFound
	}
	return t.transactionRepo.GetTransitiveSum(transactionID)
}
