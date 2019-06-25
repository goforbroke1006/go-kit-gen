install:
	GOBIN=${GOPATH}/bin go install ./cmd/go-kit-gen

test:
	go test ./...