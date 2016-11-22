install:
	go install -v

build:
	go build -v ./...

lint:
	golint ./...
	go vet ./...

test:
	go test -v ./... --cover

deps:
	go get github.com/nats-io/nats
	go get github.com/ernestio/ernest-config-client
	go get github.com/ernestio/ernestaws
	go get github.com/aws/aws-sdk-go/aws
	go get github.com/satori/go.uuid

dev-deps: deps
	go get github.com/golang/lint/golint
	go get github.com/smartystreets/goconvey/convey

clean:
	go clean

dist-clean:
	rm -rf pkg src bin
