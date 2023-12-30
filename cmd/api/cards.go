package main

import (
	"database/sql"
	"errors"
	"net/http"

	"ratrace.darieldejesus.com/internal/data"
	"ratrace.darieldejesus.com/internal/validator"
)

func (app *application) createCardHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Category       string `json:"category"`
		Title          string `json:"title"`
		Body           string `json:"body"`
		Type           string `json:"type"`
		Cost           int    `json:"cost"`
		DownPayment    int    `json:"down_payment"`
		Mortgage       int    `json:"mortgage"`
		CashFlow       int    `json:"cash_flow"`
		TradeRangeDown int    `json:"trade_range_down"`
		TradeRangeUp   int    `json:"trade_range_up"`
		Inflation      int    `json:"inflation"`
		Quantity       int    `json:"quantity"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	card := &data.Card{
		Category:       input.Category,
		Title:          input.Title,
		Body:           input.Body,
		Type:           input.Type,
		Cost:           input.Cost,
		DownPayment:    input.DownPayment,
		Mortgage:       input.Mortgage,
		CashFlow:       input.CashFlow,
		TradeRangeDown: input.TradeRangeDown,
		TradeRangeUp:   input.TradeRangeUp,
		Inflation:      input.Inflation,
		Quantity:       input.Quantity,
		Active:         true,
	}

	v := validator.New()
	if data.ValidateCard(v, card); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Cards.Insert(card)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"card": card}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showCardHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	card, err := app.models.Cards.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"card": card}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
