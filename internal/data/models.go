package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("")
)

type Models struct {
	Parties     PartyModel
	Users       UserModel
	Tokens      TokenModel
	Permissions PermissionModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Parties:     PartyModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
