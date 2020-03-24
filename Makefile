.PHONY: build tidy swag pb

build: tidy
	@go build -race -o ./bin/go-app.exe -gcflags "all=-N -l" ./main.go

tidy:
	@go mod tidy

pb:
	@protoc --proto_path=. --go_out=. ./api/hello.proto
	#./scripts/pb.bat

swag:
	swag i --dir ./

