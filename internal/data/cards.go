package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"ratrace.darieldejesus.com/internal/validator"
)

type Card struct {
	ID             int64     `json:"id"`
	Category       string    `json:"category"`
	Title          string    `json:"title"`
	Body           string    `json:"body"`
	Type           string    `json:"type"`
	Cost           int       `json:"cost"`
	DownPayment    int       `json:"down_payment"`
	Mortgage       int       `json:"mortgage"`
	CashFlow       int       `json:"cash_flow"`
	TradeRangeDown int       `json:"trade_range_down"`
	TradeRangeUp   int       `json:"trade_range_up"`
	Inflation      int       `json:"inflation"`
	Quantity       int       `json:"quantity"`
	Active         bool      `json:"active"`
	CreatedAt      time.Time `json:"-"`
}

func ValidateCard(v *validator.Validator, c *Card) {
	allowedCategories := []string{"investment", "small_investment", "expense", "market"}
	v.Check(validator.In(c.Category, allowedCategories...), "category", fmt.Sprintf("must be provided a valid category (%s)", strings.Join(allowedCategories, ", ")))
	v.Check(validator.NotBlank(c.Title), "title", "must be provided")
	v.Check(validator.NotBlank(c.Body), "body", "must be provided")
	v.Check(validator.NotBlank(c.Type), "type", "must be provided")
	v.Check(c.Cost > 0, "cost", "must be greather than zero")
	v.Check(c.Quantity > 0, "quantity", "must be greather than zero")
}

type CardModel struct {
	DB *sql.DB
}

func (m CardModel) Insert(c *Card) error {
	stmt := `INSERT INTO cards (category, title, body, type, cost, down_payment, mortgage, cash_flow, trade_range_down, trade_range_up, inflation, quantity, active)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := m.DB.ExecContext(
		ctx,
		stmt,
		c.Category,
		c.Title,
		c.Body,
		c.Type,
		c.Cost,
		c.DownPayment,
		c.Mortgage,
		c.CashFlow,
		c.TradeRangeDown,
		c.TradeRangeUp,
		c.Inflation,
		c.Quantity,
		c.Active,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	c.ID = id
	return nil
}

func (m CardModel) Get(id int64) (*Card, error) {
	stmt := `
	SELECT id, category, title, body, type, cost, down_payment, mortgage, cash_flow, trade_range_down, trade_range_up, inflation, quantity, active, created_at
	FROM cards
	WHERE id = ? AND active = true`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var card Card
	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(
		&card.ID,
		&card.Category,
		&card.Title,
		&card.Body,
		&card.Type,
		&card.Cost,
		&card.DownPayment,
		&card.Mortgage,
		&card.CashFlow,
		&card.TradeRangeDown,
		&card.TradeRangeUp,
		&card.Inflation,
		&card.Quantity,
		&card.Active,
		&card.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &card, nil
}
