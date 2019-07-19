SERVICE_NAME=go-kit-gen

build:
	go build -o build/release/${SERVICE_NAME} ./cmd/${SERVICE_NAME}

deps:
	go get github.com/go-kit/kit/endpoint
	go get github.com/go-kit/kit/transport
	dep ensure -v

install:
	GOBIN=${GOPATH}/bin go install ./cmd/${SERVICE_NAME}

.PHONY: test
test:
#	protoc \
#		--proto_path=$(GOPATH)/src/github.com/gogo/protobuf/protobuf/ \
#		--proto_path=. \
#		--go_out=plugins=grpc:. \
#		./test/testdata/some-awesome-hub.proto
	go test ./...

test-cover:
	go test -cover ./...