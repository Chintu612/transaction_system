package repositories

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"transaction_system/app/lib/db"
	"transaction_system/app/models"

	"gorm.io/gorm"
)

//go:generate mockgen -source=./transaction.go -destination=mock_repositories/mock_transaction.go -package=mock_repositories

const (
	// DuplicateKeyViolationCode postgreSQL error code 23505 corresponds to a unique_violation (duplicate key)
	DuplicateKeyViolationCode = "23505"
)

var ErrTransactionAlreadyExist = errors.New("transaction with the same ID already exists")

type TransactionRepositoryI interface {
	Create(transaction *models.Transaction) error
	GetByID(transactionID uint) (*models.Transaction, error)
	GetByType(transactionType string) ([]models.Transaction, error)
	GetTransitiveSum(transactionID uint) (float64, error)
}

type transactionRepository struct {
	Db *gorm.DB
}

func NewTransactionRepository() TransactionRepositoryI {
	return &transactionRepository{Db: db.Get()}
}

// Create inserts a new transaction into the database.
func (t *transactionRepository) Create(transaction *models.Transaction) error {
	if err := t.Db.Create(transaction).Error; err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == DuplicateKeyViolationCode {
			return ErrTransactionAlreadyExist
		}
		return err
	}
	return nil
}

// GetByID retrieves a transaction by its ID from the database.
func (t *transactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	result := t.Db.Where("id = ?", id).First(&transaction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// No record found
			return nil, nil
		}
		// Other error occurred
		return nil, result.Error
	}

	return &transaction, nil
}

// GetByType retrieves transactions by type from the database.
func (t *transactionRepository) GetByType(transactionType string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := t.Db.Where("type = ?", transactionType).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

// GetTransitiveSum retrieves the sum of all transactions transitively linked by their parent_id to a given transaction ID.
func (t *transactionRepository) GetTransitiveSum(transactionID uint) (float64, error) {
	var totalAmount float64
	query := fmt.Sprintf(`
		WITH RECURSIVE TransactionsCTE AS (
			SELECT id, amount
			FROM transactions
			WHERE id = %d
	
			UNION ALL
	
			SELECT t.id, t.amount
			FROM transactions t
			JOIN TransactionsCTE ON t.parent_id = TransactionsCTE.id
		)
		SELECT COALESCE(SUM(amount), 0) AS total_amount
		FROM TransactionsCTE where id != %d;
	`, transactionID, transactionID)

	err := t.Db.Raw(query).Row().Scan(&totalAmount)
	if err != nil {
		return 0, err
	}

	return totalAmount, nil
}
