package accounts

import "time"

type (
	// Member model structure
	Member struct {
		ID        string    `json:"id"`
		AccountID string    `json:"account_id"`
		UserID    string    `json:"user_id"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	}
)
