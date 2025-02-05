.PHONY: generate build test

generate: clean
	mkdir -p ./api/pb
	protoc --experimental_allow_proto3_optional \
		--go_out=./api/pb \
		--go-grpc_out=./api/pb \
		-I./proto \
		proto/*.proto

clean:
	rm -rf api/pb/*

build:
	go build -o bin/server cmd/main.go

test:
	go test ./...

run:
	go run cmd/main.go