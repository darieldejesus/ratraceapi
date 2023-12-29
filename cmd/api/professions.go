package main

import (
	"errors"
	"net/http"

	"ratrace.darieldejesus.com/internal/data"
	"ratrace.darieldejesus.com/internal/validator"
)

func (app *application) createProfessionHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name          string `json:"name"`
		Salary        int    `json:"salary"`
		Savings       int    `json:"savings"`
		Taxes         int    `json:"taxes"`
		Mortgage      int    `json:"mortgage"`
		SchoolLoan    int    `json:"school_loan"`
		CarLoan       int    `json:"car_loan"`
		CreditCard    int    `json:"credit_card"`
		OtherExpenses int    `json:"other_expenses"`
		BankLoan      int    `json:"bank_loan"`
		ChildExpense  int    `json:"child_expense"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	profession := &data.Profession{
		Name:          input.Name,
		Salary:        input.Salary,
		Savings:       input.Savings,
		Taxes:         input.Taxes,
		Mortgage:      input.Mortgage,
		SchoolLoan:    input.SchoolLoan,
		CarLoan:       input.CarLoan,
		CreditCard:    input.CreditCard,
		OtherExpenses: input.OtherExpenses,
		BankLoan:      input.BankLoan,
		ChildExpense:  input.ChildExpense,
		Active:        true,
	}

	v := validator.New()
	if data.ValidateProfession(v, profession); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Professions.Insert(profession)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"profession": profession}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showProfessionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	profession, err := app.models.Professions.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"profession": profession}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
