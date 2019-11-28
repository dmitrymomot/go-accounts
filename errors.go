package accounts

import "errors"

// Predefined accounts package errors
var (
	ErrNameMissed        = errors.New("name is missed")
	ErrOwnerIDMissed     = errors.New("owner id is missed")
	ErrNotExistedAccount = errors.New("could not update not existed account")
)
