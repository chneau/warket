.SILENT:
.ONESHELL:
.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.PHONY: run deps build clean exec test

name=$(shell basename $(CURDIR))

run: build exec clean

exec:
	./bin/${name}

build:
	CGO_ENABLED=0 go build -o bin/${name} -ldflags '-s -w -extldflags "-static"'

test:
	go test --count=1 ./pkg/client/...

dist:
	gox -verbose -ldflags '-s -w -extldflags "-static"' -osarch="linux/amd64 windows/amd64" -output "dist/{{.OS}}_{{.Arch}}_{{.Dir}}"

release:
	ghr --delete --replace --prerelease --debug pre-release dist/

dev-dist:
	go get github.com/mitchellh/gox
	go get github.com/tcnksm/ghr

clean:
	rm -rf bin dist

# deps:
# 	GO111MODULE=on go mod vendor -v

deps:
	govendor init
	govendor add +e
	govendor update +v

dev:
	go get -u -v github.com/kardianos/govendor
