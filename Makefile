.PHONY: build

build:
	go build -race -o go-app.exe -gcflags "all=-N -l" ./cmd/main.go
