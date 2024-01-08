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
	Boards      BoardModel
	Cards       CardModel
	Parties     PartyModel
	Permissions PermissionModel
	Professions ProfessionModel
	Tokens      TokenModel
	Users       UserModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Boards:      BoardModel{DB: db},
		Cards:       CardModel{DB: db},
		Parties:     PartyModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Professions: ProfessionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
