.SILENT:
.ONESHELL:
.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.PHONY: run deps build clean exec test

NAME=$(shell basename $(CURDIR))

run: build exec clean

exec:
	./bin/${NAME}

build:
	CGO_ENABLED=0 go build -trimpath -o bin/${NAME} -ldflags '-s -w -extldflags "-static"'

test:
	go test --count=1 ./client/...

dist:
	gox -verbose -ldflags '-s -w -extldflags "-static"' -osarch="linux/amd64 windows/amd64" -output "dist/{{.OS}}_{{.Arch}}_{{.Dir}}"

release:
	ghr --delete --replace --prerelease --debug pre-release dist/

dev-dist:
	go get github.com/mitchellh/gox
	go get github.com/tcnksm/ghr

clean:
	rm -rf bin dist

deps:
	rm -f go.mod go.sum
	go mod init || true
	go get -u ./...
	go mod tidy

watch_s:
	nodemon -e go -x "go run . s -g 0 || false"
