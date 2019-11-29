package accounts

type (
	// Account model structure
	Account struct {
		ID        string `db:"id" json:"id"`
		Name      string `db:"name" json:"name"`
		Disabled  bool   `db:"disabled" json:"disabled"`
		CreatedAt int64  `db:"created_at" json:"created_at"`
		UpdatedAt *int64 `db:"updated_at" json:"updated_at,omitempty,"`
	}
)
