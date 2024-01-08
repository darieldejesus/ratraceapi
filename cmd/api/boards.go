package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"ratrace.darieldejesus.com/internal/data"
)

func (app *application) createBoardHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name         string   `json:"name"`
		Type         string   `json:"type"`
		InnerActions []string `json:"inner_actions"`
		OuterActions []string `json:"outer_actions"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	board := &data.Board{
		Name:         input.Name,
		Type:         input.Type,
		InnerActions: input.InnerActions,
		OuterActions: input.OuterActions,
		Active:       true,
	}

	// v := validator.New()
	// if data.ValidateCard(v, board); !v.Valid() {
	// 	app.failedValidationResponse(w, r, v.Errors)
	// 	return
	// }

	err = app.models.Boards.Insert(board)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"board": board}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showBoardHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	board, err := app.models.Boards.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	fmt.Println(board)
}
