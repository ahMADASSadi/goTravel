# Makefile

APP_NAME=main

.PHONY: swag run dev sync

swag:
	swag init --generalInfo cmd/main.go --parseDependency --parseInternal

run: swag
	go run /cmd/main.go

dev:
	air

sync:
	go get -u