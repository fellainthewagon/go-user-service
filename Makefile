format:
	go fmt ./...

build:
	go build -o ./build/app cmd/app.go

run:
	go run cmd/app.go