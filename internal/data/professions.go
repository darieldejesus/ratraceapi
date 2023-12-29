package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ratrace.darieldejesus.com/internal/validator"
)

type Profession struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Salary        int       `json:"salary"`
	Savings       int       `json:"savings"`
	Taxes         int       `json:"taxes"`
	Mortgage      int       `json:"mortgage"`
	SchoolLoan    int       `json:"school_loan"`
	CarLoan       int       `json:"car_loan"`
	CreditCard    int       `json:"credit_card"`
	OtherExpenses int       `json:"other_expenses"`
	BankLoan      int       `json:"bank_loan"`
	ChildExpense  int       `json:"child_expense"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"-"`
}

func ValidateProfession(v *validator.Validator, p *Profession) {
	v.Check(validator.NotBlank(p.Name), "name", "must be provided")
	v.Check(validator.MaxChars(p.Name, 255), "title", "must not be more than 255 bytes long")
	v.Check(p.Salary >= 0, "salary", "must not be a negative number")
	v.Check(p.Savings >= 0, "savings", "must not be a negative number")
	v.Check(p.Taxes >= 0, "taxes", "must not be a negative number")
	v.Check(p.Mortgage >= 0, "mortgage", "must not be a negative number")
	v.Check(p.SchoolLoan >= 0, "school_loan", "must not be a negative number")
	v.Check(p.CarLoan >= 0, "car_loan", "must not be a negative number")
	v.Check(p.CreditCard >= 0, "credit_card", "must not be a negative number")
	v.Check(p.OtherExpenses >= 0, "other_expenses", "must not be a negative number")
	v.Check(p.BankLoan >= 0, "bank_loan", "must not be a negative number")
	v.Check(p.ChildExpense >= 0, "child_expense", "must not be a negative number")
}

type ProfessionModel struct {
	DB *sql.DB
}

func (m ProfessionModel) Insert(p *Profession) error {
	stmt := `INSERT INTO professions (name, salary, savings, taxes, mortgage, school_loan, car_loan, credit_card, other_expenses, bank_loan, child_expense, active)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := m.DB.ExecContext(
		ctx,
		stmt,
		p.Name,
		p.Salary,
		p.Savings,
		p.Taxes,
		p.Mortgage,
		p.SchoolLoan,
		p.CarLoan,
		p.CreditCard,
		p.OtherExpenses,
		p.BankLoan,
		p.ChildExpense,
		p.Active,
	)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = id
	return nil
}

func (m ProfessionModel) Get(id int64) (*Profession, error) {
	stmt := `
	SELECT id, name, salary, savings, taxes, mortgage, school_loan, car_loan, credit_card, other_expenses, bank_loan, child_expense, active, created_at
	FROM professions
	WHERE id = ? AND active = true`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var profession Profession
	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(
		&profession.ID,
		&profession.Name,
		&profession.Salary,
		&profession.Savings,
		&profession.Taxes,
		&profession.Mortgage,
		&profession.SchoolLoan,
		&profession.CarLoan,
		&profession.CreditCard,
		&profession.OtherExpenses,
		&profession.BankLoan,
		&profession.ChildExpense,
		&profession.Active,
		&profession.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &profession, nil
}
