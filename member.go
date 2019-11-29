package accounts

type (
	// Member model structure
	Member struct {
		ID        string `json:"id" db:"id"`
		AccountID string `json:"account_id" db:"account_id"`
		UserID    string `json:"user_id" db:"user_id"`
		Role      string `json:"role" db:"role"`
		CreatedAt int64  `json:"created_at" db:"created_at"`
	}
)
