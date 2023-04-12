GOPATH:=$(shell go env GOPATH)
MOBILESERVICE:= mobile-fishing
ADMINSERVICE:= admin-fishing

.PHONY: proto
proto:
	protoc -I=../../../ --proto_path=${GOPATH}/src:. --realmicro_out=. --go_out=. proto/mobile/${MOBILESERVICE}.proto
	protoc  -I=../../../ --proto_path=${GOPATH}/src:. --realmicro_out=. --go_out=. proto/admin/${ADMINSERVICE}.proto

.PHONY: adminProto
adminProto:
	protoc  -I=../../../ --proto_path=${GOPATH}/src:. --realmicro_out=. --go_out=. proto/admin/${ADMINSERVICE}.proto

.PHONY: mobileProto
mobileProto:
	protoc  -I=../../../ --proto_path=${GOPATH}/src:. --realmicro_out=. --go_out=. proto/mobile/${MOBILESERVICE}.proto



.PHONY: build
build: proto
	go build  main.go

.PHONY: test
test:
	go test -v ./... -cover -gcflags=-l
