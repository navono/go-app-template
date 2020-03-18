.PHONY: build swag

build:
	go build -race -o ./bin/go-app.exe -gcflags "all=-N -l" ./main.go

swag:
	swag i --dir ./

