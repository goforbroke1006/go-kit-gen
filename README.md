# go-kit-gen

Tool for generating gokit-based microservice project structure with proto-3 file.

**The project is not finished yet!**

### CLI

```bash
go-kit-gen \
    --working-dir=/home/user/go/src/github.com/goforbroke1006/test-svc \
    --proto-path=/home/user/go/src/github.com/goforbroke1006/test-svc \
    --proto-file=pb/api/v1/test-service.proto \
    --service-name=SomeAwesomeHub
```

### Project structure

* pkg
    * endpoint
    * service
    * transport
    * model
