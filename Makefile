VERSION := v0.1.2

APP := DeNet-node
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)

DIR = builds/

HOSTOS = $(APP)-$(VERSION)-$(OS)-$(ARCH)
WINAMD64 = $(APP)-$(VERSION)-windows-amd64.exe
WIN386 = $(APP)-$(VERSION)-windows-i386.exe

build:
	go build -ldflags "-s -w" -o $(DIR)$(HOSTOS)
	GOOS=windows CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -ldflags "-s -w" -o $(DIR)$(WINAMD64)
	GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ go build -ldflags "-s -w" -o $(DIR)$(WIN386)

hosted_os:
	go build -ldflags "-s -w" -o $(DIR)$(HOSTOS)

update_config:
	rm -rf docs
	swag init -g server/server.go

gen:
	protoc -I=./proto --go_out=./proto --go-grpc_out=./proto upload.proto

clean:
	rm ./pb/*.go