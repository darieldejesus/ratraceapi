package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ratrace.darieldejesus.com/internal/validator"
)

type Party struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Active       bool      `json:"active"`
	Participants []int64   `json:"participants"`
	CreatedAt    time.Time `json:"-"`
}

type PartyModel struct {
	DB *sql.DB
}

func ValidateParty(v *validator.Validator, party *Party) {
	v.Check(validator.NotBlank(party.Title), "title", "must be provided")
	v.Check(validator.MaxChars(party.Title, 50), "title", "must not be more than 500 bytes long")

	// v.Check(validator.NotEmptyList(party.Participants), "participants", "must be provided")
	// v.Check(validator.Unique(party.Participants), "participants", "must not contain duplicated values")
}

func (m PartyModel) Insert(party *Party) error {
	stmt := `INSERT INTO parties (title, active, createdAt)
	VALUES (?, ?, UTC_TIMESTAMP())`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, stmt, party.Title, party.Active)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	party.ID = id
	return nil
}

func (m PartyModel) Get(id int64) (*Party, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	stmt := `SELECT sleep(20), id, title, active, createdAt 
	FROM parties WHERE id = ? AND active = true`

	var party Party

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(
		&[]byte{},
		&party.ID,
		&party.Title,
		&party.Active,
		&party.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &party, nil
}

func (m PartyModel) Update(party *Party) error {
	stmt := `UPDATE parties SET title = ?
	WHERE id = ? AND active = true`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt, party.Title, party.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m PartyModel) Delete(id int64) error {
	stmt := `UPDATE parties SET active = false WHERE id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}
