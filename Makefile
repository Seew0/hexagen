PORT ?= 8080

run:
	go run ./main.go -i

build:
	go build -o bin/app ./main.go