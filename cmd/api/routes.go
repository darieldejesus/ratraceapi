package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)

	// Parties
	router.HandlerFunc(http.MethodGet, "/v1/parties", app.requirePermission("parties:read", app.listPartiesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/parties", app.requirePermission("parties:write", app.createPartyHandler))
	router.HandlerFunc(http.MethodGet, "/v1/parties/:id", app.requirePermission("parties:read", app.showPartyHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/parties/:id", app.requirePermission("parties:write", app.updatePartyHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/parties/:id", app.requirePermission("parties:write", app.deletePartyHandler))

	// Professions
	router.HandlerFunc(http.MethodPost, "/v1/professions", app.createProfessionHandler)
	router.HandlerFunc(http.MethodGet, "/v1/professions/:id", app.showProfessionHandler)

	// Cards
	router.HandlerFunc(http.MethodPost, "/v1/cards", app.requirePermission("cards:write", app.createCardHandler))
	router.HandlerFunc(http.MethodGet, "/v1/cards/:id", app.requirePermission("cards:read", app.showCardHandler))

	// Tokens
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authenticate", app.createAuthTokenHandler)

	// Users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activate", app.activateUserHandler)

	// Debug
	router.HandlerFunc(http.MethodGet, "/system/status", app.requirePermission("system:read", expvar.Handler().ServeHTTP))

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
