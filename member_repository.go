package accounts

import "github.com/jmoiron/sqlx"

// Predefined order types
const (
	// ORDER BY created_at ASC
	MemberOrderAsc MemberOrder = iota + 1
	// ORDER BY created_at DESC
	MemberOrderDesc
)

var memberOrderQuery = map[MemberOrder]string{
	MemberOrderAsc:  "ORDER BY created_at ASC",
	MemberOrderDesc: "ORDER BY created_at DESC",
}

type (
	// MemberOrder type
	MemberOrder int

	// MembRepository struct
	MembRepository struct {
		db        *sqlx.DB
		tableName string
	}
)

// NewMemberRepository factory
func NewMemberRepository(db *sqlx.DB, tableName string) *MembRepository {
	return &MembRepository{db: db, tableName: tableName}
}
