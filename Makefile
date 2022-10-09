run:
	go run cmd/main.go

build:
	go build -ldflags="-s -w" -trimpath -o bin/main cmd/main.go