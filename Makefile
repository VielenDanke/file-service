GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/unistack-org/protoc-gen-micro
	
.PHONY: build
build:
	go build -o user-service.exe *.go

.PHONY: proto
proto:
	protoc -I. \
        -IC:/Users/viele/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-IC:/Users/viele/go/src/github.com/grpc-ecosystem/grpc-gateway \
        --openapiv2_out=disable_default_errors=true,allow_merge=true:. --go_out=:. --micro_out=components="micro|http|grpc|gorilla":. proto/*.proto

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t file-service:latest