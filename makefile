# Makefile

APP_NAME=main

default: dev

.PHONY: swag run dev sync

swag:
	swag init -g cmd/main.go

run: swag
	go run cmd/main.go

dev: swag
	air

sync:
	go get -u

doc:
	godoc -http=:6060
