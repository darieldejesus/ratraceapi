package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Parties PartyModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Parties: PartyModel{DB: db},
	}
}
