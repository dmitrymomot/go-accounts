package accounts

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type (
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

// GetByID ...
func (r *MembRepository) GetByID(id string) (*Member, error) {
	q := "SELECT * FROM %s WHERE id = ? LIMIT 1"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	m := &Member{}
	if err := r.db.Get(m, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFoundMember
		}
		return nil, errors.Wrap(err, "get member by id")
	}
	return m, nil
}

// GetList ...
func (r *MembRepository) GetList(c ...Condition) ([]*Member, error) {
	q := "SELECT * FROM %s "
	q = fmt.Sprintf(q, r.tableName)
	sq, params := conditionsToQuery(c...)
	q = q + sq
	q = r.db.Rebind(q)
	al := make([]*Member, 0)
	if err := r.db.Select(&al, q, params...); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFoundMember
		}
		return nil, errors.Wrap(err, "get members list")
	}
	return al, nil
}

// Insert ...
func (r *MembRepository) Insert(m *Member) error {
	q := "INSERT INTO %s (`id`, `account_id`, `user_id`, `role`, `created_at`) VALUES (?, ?, ?, ?, ?);"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(q, m.ID, m.AccountID, m.UserID, m.Role, m.CreatedAt); err != nil {
		return errors.Wrap(err, "store member")
	}
	return nil
}

// Update ...
func (r *MembRepository) Update(m *Member) error {
	q := "UPDATE %s SET `role`=? WHERE id=?;"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(q, m.Role, m.ID); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFoundMember
		}
		return errors.Wrap(err, "update member")
	}
	return nil
}

// Delete ...
func (r *MembRepository) Delete(id string) error {
	q := "DELETE FROM %s WHERE id=?;"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(q, id); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFoundMember
		}
		return errors.Wrap(err, "delete member")
	}
	return nil
}

// DeleteByAccountID ...
func (r *MembRepository) DeleteByAccountID(aid string) error {
	q := "DELETE FROM %s WHERE account_id=?;"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(q, aid); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFoundMember
		}
		return errors.Wrap(err, "delete member by account id")
	}
	return nil
}

// DeleteByUserID ...
func (r *MembRepository) DeleteByUserID(uid string) error {
	q := "DELETE FROM %s WHERE user_id=?;"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(q, uid); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFoundMember
		}
		return errors.Wrap(err, "delete member by user id")
	}
	return nil
}

// HasRole checks whether the member has the given role in the account or not
func (r *MembRepository) HasRole(aid, uid, role string) error {
	q := "SELECT * FROM %s WHERE account_id = ? AND user_id = ? AND role = ? LIMIT 1"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	m := &Member{}
	if err := r.db.Get(m, q, aid, uid, role); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFoundMember
		}
		return errors.Wrap(err, "has member role")
	}
	return nil
}
