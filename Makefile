BINARY_NAME=tfpvc

build-all:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-osx-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-osx-arm64 main.go

install-linux:
	cp bin/${BINARY_NAME}-linux-amd64 /usr/local/bin/tfpvc

install-mac-intel:
	cp bin/${BINARY_NAME}-osx-amd64 /usr/local/bin/tfpvc

install-mac-applesilicon:
	cp bin/$(BINARY_NAME)-osx-arm64 /usr/local/bin/tfpvc

all: build-all