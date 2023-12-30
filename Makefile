include .envrc

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

###################################
# DEVELOPMENT
###################################

## run: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api \
		-port=${RR_PORT} \
		-db-dsn="${RR_DB_DSN}" \
		-db-max-open-conns=${RR_DB_MAX_OPEN_CONNS} \
		-db-max-idle-conns=${RR_DB_MAX_IDLE_CONNS} \
		-db-max-idle-time=${RR_DB_MAX_IDLE_TIME} \
		-limiter-enabled=${RR_LIMITER_ENABLED} \
		-limiter-rps=${RR_LIMITER_RPS} \
		-limiter-burst=${RR_LIMITER_BURST} \
		-smtp-host=${RR_SMTP_HOST} \
		-smtp-port=${RR_SMTP_PORT} \
		-smtp-username=${RR_SMTP_USERNAME} \
		-smtp-password=${RR_SMTP_PASSWORD} \
		-smtp-sender=${RR_SMTP_SENDER} \
		-cors-trusted-origins=${RR_CORS_TRUSTED_ORIGINS}

## migrations/new name=$1: create a new database migration
.PHONY: migrations/new
migrations/new:
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## migrations/up: apply all up database migrations
.PHONY: migrations/up
migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database="mysql://${RR_DB_DSN}" up

###################################
# QUALITY CONTROL
###################################
.PHONY: audit
audit:
	@echo 'Verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting source code...'
	go fmt ./...
	@echo 'Vetting source code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...
