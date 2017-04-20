APP=smartsc
VERSION=latest

all: deps build

clean:
	@echo "--> cleaning..."
	@rm -rf build
	@rm -rf vendor
	@go clean ./...

prereq:
	@mkdir -p build/{bin,tar}
	@go get -u github.com/Masterminds/glide

deps: prereq
	@glide install

build: prereq
	@echo '--> building...'
	@go fmt ./...
	go build -o build/bin/${APP} ./cmd

package:
	@echo '--> packaging...'
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -o build/bin/${APP} ./cmd
	@docker build -t reg.skycloud.com:5000/firmament/smartsc:${VERSION} .

deploy: package
	@echo '--> deploying...'
	@docker push reg.skycloud.com:5000/firmament/smartsc:${VERSION}
