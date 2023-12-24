package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/parties", app.createPartyHandler)
	router.HandlerFunc(http.MethodGet, "/v1/parties/:id", app.showPartyHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/parties/:id", app.updatePartyHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/parties/:id", app.deletePartyHandler)

	return router
}
