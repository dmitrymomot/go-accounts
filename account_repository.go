package accounts

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type (
	// AccRepository struct
	AccRepository struct {
		db        *sqlx.DB
		tableName string
	}
)

// NewAccountRepository factory
func NewAccountRepository(db *sql.DB, driverName, tableName string) *AccRepository {
	return &AccRepository{
		db:        sqlx.NewDb(db, driverName),
		tableName: tableName,
	}
}

// GetByID ...
func (r *AccRepository) GetByID(id string) (*Account, error) {
	q := "SELECT * FROM %s WHERE id = ? LIMIT 1"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	a := &Account{}
	if err := r.db.Get(a, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFoundAccount
		}
		return nil, errors.Wrap(err, "get account by id")
	}
	return a, nil
}

// GetList ...
func (r *AccRepository) GetList(c ...Condition) ([]*Account, error) {
	q := "SELECT * FROM %s "
	q = fmt.Sprintf(q, r.tableName)
	sq, params := conditionsToQuery(c...)
	q = q + sq
	q = r.db.Rebind(q)
	al := make([]*Account, 0)
	if err := r.db.Select(&al, q, params...); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFoundAccount
		}
		return nil, errors.Wrap(err, "get accounts list")
	}
	return al, nil
}

// Insert ...
func (r *AccRepository) Insert(a *Account) error {
	q := "INSERT INTO %s (`id`, `name`, `disabled`, `created_at`) VALUES (?, ?, ?, ?);"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(q, a.ID, a.Name, a.Disabled, a.CreatedAt); err != nil {
		return errors.Wrap(err, "store account")
	}
	return nil
}

// Update ...
func (r *AccRepository) Update(a *Account) error {
	q := "UPDATE %s SET `name`=?, `disabled`=?, `updated_at`=? WHERE id=?;"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(q, a.Name, a.Disabled, a.UpdatedAt, a.ID); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFoundAccount
		}
		return errors.Wrap(err, "update account")
	}
	return nil
}

// Delete ...
func (r *AccRepository) Delete(id string) error {
	q := "DELETE FROM %s WHERE id=?;"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(q, id); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFoundAccount
		}
		return errors.Wrap(err, "delete account")
	}
	return nil
}
