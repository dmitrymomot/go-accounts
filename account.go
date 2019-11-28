package accounts

import "time"

type (
	// Account model structure
	Account struct {
		ID        string     `json:"id"`
		Name      string     `json:"name"`
		Disabled  bool       `json:"disabled"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at,omitempty,"`
	}
)
