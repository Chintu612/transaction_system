package models

// Transaction represents the transactions table schema.
type Transaction struct {
	Id       uint    `json:"id" gorm:"primarykey"`
	Amount   float64 `json:"amount" validate:"notblank"`
	Type     string  `json:"type" validate:"notblank" gorm:"varchar(50)"`
	ParentID *uint   `json:"parent_id"`
}

func (Transaction) TableName() string {
	return "transactions"
}
