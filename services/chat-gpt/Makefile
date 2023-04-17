GOPATH:=$(shell go env GOPATH)
MOBILESERVICE:= mobile-chat-gpt
ADMINSERVICE:= admin-fishing

.PHONY: proto
proto:
	protoc -I=../../../ --proto_path=${GOPATH}/src:. --realmicro_out=. --go_out=. proto/mobile/${MOBILESERVICE}.proto


.PHONY: build
build: proto
	go build  main.go

.PHONY: test
test:
	go test -v ./... -cover -gcflags=-l
