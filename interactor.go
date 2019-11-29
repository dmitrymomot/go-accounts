package accounts

import (
	"time"

	"github.com/google/uuid"
)

type (
	// Interactor struct
	Interactor struct {
		accRepo  AccountRepository
		membRepo MemberRepository
	}

	// AccountRepository interface
	AccountRepository interface {
		GetByID(id string) (*Account, error)
		GetList(c ...Condition) ([]*Account, error)
		Insert(a *Account) error
		Update(a *Account) error
		Delete(id string) error
	}

	// MemberRepository interface
	MemberRepository interface {
		GetByID(id string) (*Member, error)
		HasRole(aid, uid, role string) error
		GetList(c ...Condition) ([]*Member, error)
		Insert(m *Member) error
		Update(m *Member) error
		Delete(id string) error
		DeleteByAccountID(aid string) error
		DeleteByUserID(uid string) error
	}

	// AccountWithRole is account model with role of user for each account
	AccountWithRole struct {
		Account
		Role string
	}
)

// NewInteractor factory
func NewInteractor(ar AccountRepository, mr MemberRepository) *Interactor {
	return &Interactor{
		accRepo:  ar,
		membRepo: mr,
	}
}

// GetByID fetch account by id
func (i *Interactor) GetByID(id string) (*Account, error) {
	return i.accRepo.GetByID(id)
}

// GetList of accounts
func (i *Interactor) GetList(c ...Condition) ([]*Account, error) {
	return i.accRepo.GetList(c...)
}

// Insert a new account and add user as owner
func (i *Interactor) Insert(a *Account, uid string) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if a.CreatedAt == 0 {
		a.CreatedAt = time.Now().Unix()
	}
	if a.Name == "" {
		return ErrNameMissed
	}
	if uid == "" {
		return ErrOwnerIDMissed
	}
	if err := i.accRepo.Insert(a); err != nil {
		return err
	}
	if err := i.membRepo.Insert(&Member{
		ID:        uuid.New().String(),
		AccountID: a.ID,
		UserID:    uid,
		Role:      RoleOwner,
		CreatedAt: time.Now().Unix(),
	}); err != nil {
		i.accRepo.Delete(a.ID)
		return err
	}
	return nil
}

// Update existed account
func (i *Interactor) Update(a *Account) error {
	if a.ID == "" {
		return ErrNotExistedAccount
	}
	if a.UpdatedAt == nil {
		t := time.Now().Unix()
		a.UpdatedAt = &t
	}
	if a.Name == "" {
		return ErrNameMissed
	}
	return i.accRepo.Update(a)
}

// Delete account by id with all related members
func (i *Interactor) Delete(id string) error {
	if err := i.accRepo.Delete(id); err != nil {
		return err
	}
	if err := i.membRepo.DeleteByAccountID(id); err != nil {
		return err
	}
	return nil
}

// AddMember member to account
func (i *Interactor) AddMember(aid, uid, role string) error {
	return i.membRepo.Insert(&Member{
		ID:        uuid.New().String(),
		AccountID: aid,
		UserID:    uid,
		Role:      role,
		CreatedAt: time.Now().Unix(),
	})
}

// DeleteMemberByID deletes member by member id
func (i *Interactor) DeleteMemberByID(mid string) error {
	return i.membRepo.Delete(mid)
}

// DeleteMembersByAccountID deletes members by account id
func (i *Interactor) DeleteMembersByAccountID(aid string) error {
	return i.membRepo.DeleteByAccountID(aid)
}

// DeleteMembersByUserID deletes members by user id
func (i *Interactor) DeleteMembersByUserID(uid string) error {
	return i.membRepo.DeleteByUserID(uid)
}

// GetMembersList returns members list
func (i *Interactor) GetMembersList(c ...Condition) ([]*Member, error) {
	ml, err := i.membRepo.GetList(c...)
	if err != nil {
		return nil, err
	}
	return ml, nil
}

// GetAccountsListByUserID returns accounts list by user id
func (i *Interactor) GetAccountsListByUserID(uid string, c ...Condition) ([]*Account, error) {
	if len(c) > 0 {
		c = append(c, UserID(uid))
	} else {
		c = []Condition{UserID(uid)}
	}
	ml, err := i.membRepo.GetList(c...)
	if err != nil {
		return nil, err
	}
	al, err := i.accRepo.GetList(IDs(i.getAccountsIDs(ml)))
	if err != nil {
		return nil, err
	}
	return al, nil
}

// GetAccountsListWithRoleByUserID returns accounts list with roles by user id
func (i *Interactor) GetAccountsListWithRoleByUserID(uid string, c ...Condition) ([]*AccountWithRole, error) {
	if len(c) > 0 {
		c = append(c, UserID(uid))
	} else {
		c = []Condition{UserID(uid)}
	}
	ml, err := i.membRepo.GetList(c...)
	if err != nil {
		return nil, err
	}
	al, err := i.accRepo.GetList(IDs(i.getAccountsIDs(ml)))
	if err != nil {
		return nil, err
	}

	arm := make(map[string]string)
	for _, m := range ml {
		arm[m.AccountID] = m.Role
	}

	ar := make([]*AccountWithRole, 0, len(al))
	for _, a := range al {
		if r, ok := arm[a.ID]; ok {
			item := &AccountWithRole{*a, r}
			ar = append(ar, item)
		}
	}

	return ar, nil
}

// GetUsersIDs returns slice of users ids by members list
func (i *Interactor) getUsersIDs(ml []*Member) []string {
	ids := make([]string, 0, len(ml))
	for _, m := range ml {
		ids = append(ids, m.UserID)
	}
	return ids
}

// GetAccountsIDs returns slice of accounts ids by members list
func (i *Interactor) getAccountsIDs(ml []*Member) []string {
	ids := make([]string, 0, len(ml))
	for _, m := range ml {
		ids = append(ids, m.AccountID)
	}
	return ids
}

// HasRole helper to check user role for account
func (i *Interactor) HasRole(aid, uid, role string) bool {
	if i.membRepo.HasRole(aid, uid, role) != nil {
		return false
	}
	return true
}
