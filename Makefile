.PHONY: default gen build

default: build

gen:
	go tool sqlc generate

build: gen
	GOOS=linux GOARCH=amd64 go build -o plex-query -ldflags="-s -w" -trimpath main.go

-include custom.mk

