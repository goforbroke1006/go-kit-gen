install:
	GOBIN=${GOPATH}/bin go install ./cmd/go-kit-gen

.PHONY: test
test:
	protoc --proto_path=$(GOPATH)/src/github.com/gogo/protobuf/protobuf/ \
	    --proto_path=. \
	    --go_out=plugins=grpc:. \
	    ./testdata/some-awesome-hub.proto
	go test ./...

test-cover:
	go test -cover ./...