package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Board struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	InnerActions []string  `json:"inner_actions"`
	OuterActions []string  `json:"outer_actions"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"-"`
}

type BoardModel struct {
	DB *sql.DB
}

func (m BoardModel) Insert(board *Board) error {
	stmt := `INSERT INTO boards (name, type, active)
	VALUES (?, ?, ?)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, stmt, board.Name, board.Type, board.Active)
	if err != nil {
		return err
	}

	board.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	stmt = `INSERT INTO board_spaces (board_id, type, sort, action) VALUES `

	numSpaces := len(board.InnerActions) + len(board.OuterActions)
	inserts := make([]string, 0, numSpaces)
	params := make([]any, 0, 4*numSpaces)

	for index, space := range board.InnerActions {
		inserts = append(inserts, "(?, ?, ?, ?)")
		params = append(params, board.ID, "inner", index, space)
	}

	for index, space := range board.OuterActions {
		inserts = append(inserts, "(?, ?, ?, ?)")
		params = append(params, board.ID, "outer", index, space)
	}

	queryValues := strings.Join(inserts, ",")
	stmt = stmt + queryValues

	_, err = tx.ExecContext(ctx, stmt, params...)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m BoardModel) Get(id int64) (*Board, error) {
	stmt := `SELECT id, name, type, active, created_at
	FROM boards
	WHERE id = ? AND active = true`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	board := &Board{}

	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(
		&board.ID,
		&board.Name,
		&board.Type,
		&board.Active,
		&board.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	stmt = `SELECT type, action
	FROM board_spaces
	WHERE board_id = ?
	ORDER BY sort`

	result, err := m.DB.QueryContext(ctx, stmt, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return board, nil
		default:
			return nil, err
		}
	}

	defer result.Close()

	for result.Next() {
		var spaceType string
		var spaceAction string

		result.Scan(&spaceType, &spaceAction)

		switch spaceType {
		case "inner":
			board.InnerActions = append(board.InnerActions, spaceAction)
		case "outer":
			board.OuterActions = append(board.OuterActions, spaceAction)
		}
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return board, nil
}
