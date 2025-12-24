PORT ?= 8080

run:
	go run ./cmd/app.go

build:
	go build -o bin/app ./cmd/app.go

test:
	go test ./...

setup:
	go mod tidy
