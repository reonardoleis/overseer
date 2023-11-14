run:
	go run cmd/app/main.go

build:
	go mod vendor
	go build -o bin/app cmd/app/main.go