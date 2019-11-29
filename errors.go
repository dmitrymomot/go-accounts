package accounts

import "errors"

// Predefined accounts package errors
var (
	ErrNotFoundAccount   = errors.New("account not found")
	ErrNotFoundMember    = errors.New("member not found")
	ErrNameMissed        = errors.New("name is missed")
	ErrOwnerIDMissed     = errors.New("owner id is missed")
	ErrNotExistedAccount = errors.New("could not update not existed account")
)
