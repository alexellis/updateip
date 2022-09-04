Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X github.com/alexellis/updateip/cmd.Version=$(Version) -X github.com/alexellis/updateip/cmd.GitCommit=$(GitCommit)"

SOURCE_DIRS = cmd main.go
export GO111MODULE=on

.PHONY: all
all: gofmt test build dist hash

.PHONY: build
build:
	go build

.PHONY: gofmt
gofmt:
	@test -z $(shell gofmt -l -s $(SOURCE_DIRS) ./ | tee /dev/stderr) || (echo "[WARN] Fix formatting issues with 'make gofmt'" && exit 1)

.PHONY: test
test:
	CGO_ENABLED=0 go test $(shell go list ./... | grep -v /vendor/|xargs echo) -cover

.PHONY: dist
dist:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/updateip
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/updateip-darwin
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags $(LDFLAGS) -installsuffix cgo -o bin/updateip-darwin-arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/updateip-armhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/updateip-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/updateip.exe

.PHONY: hash
hash:
	rm -rf bin/*.sha256 && ./hack/hashgen.sh
