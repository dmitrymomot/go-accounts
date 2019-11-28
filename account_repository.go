package accounts

import "github.com/jmoiron/sqlx"

// Predefined order types
const (
	// ORDER BY created_at ASC
	AccountOrderAsc AccountOrder = iota + 1
	// ORDER BY created_at DESC
	AccountOrderDesc
)

var accoutnOrderQuery = map[AccountOrder]string{
	AccountOrderAsc:  "ORDER BY created_at ASC",
	AccountOrderDesc: "ORDER BY created_at DESC",
}

type (
	// AccountOrder type
	AccountOrder int

	// AccRepository struct
	AccRepository struct {
		db        *sqlx.DB
		tableName string
	}
)

// NewAccountRepository factory
func NewAccountRepository(db *sqlx.DB, tableName string) *AccRepository {
	return &AccRepository{db: db, tableName: tableName}
}
