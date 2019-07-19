# go-kit-gen

[![Build Status](https://travis-ci.org/goforbroke1006/go-kit-gen.svg?branch=master)](https://travis-ci.org/goforbroke1006/go-kit-gen)

Tool for generating gokit-based microservice project structure with proto-3 file.

**The project is not finished yet!**

### CLI

```bash
go-kit-gen \
    --working-dir=/home/scherkesov/go/src/github.com/goforbroke1006/sport-archive-svc \
    --proto-res-file=api/pb/v1/sport-archive-svc.pb.go \
    --service-name=SportArchive \
    --transport-type=grpc
```

### Project structure

* pkg
    * endpoint
    * service
    * transport
