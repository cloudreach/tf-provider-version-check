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

release-linux:
	$(eval VERSION=$(shell git describe --tags --abbrev=0))
	sed -i "s/LOCAL/$(VERSION)/g" ./cmd/version.go
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME) main.go
	tar -C bin -czvf bin/$(BINARY_NAME).linux-amd64.tar.gz $(BINARY_NAME)
	git checkout -- ./cmd/version.go

release-mac-intel:
	$(eval VERSION=$(shell git describe --tags --abbrev=0))
	sed -i "s/LOCAL/$(VERSION)/g" ./cmd/version.go
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME) main.go
	tar -C bin -czvf bin/$(BINARY_NAME).osx-amd64.tar.gz $(BINARY_NAME)
	git checkout -- ./cmd/version.go

release-mac-applesilicon:
	$(eval VERSION=$(shell git describe --tags --abbrev=0))
	sed -i "s/LOCAL/$(VERSION)/g" ./cmd/version.go
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME) main.go
	tar -C bin -czvf bin/$(BINARY_NAME).osx-arm64.tar.gz $(BINARY_NAME)
	git checkout -- ./cmd/version.go

all: build-all