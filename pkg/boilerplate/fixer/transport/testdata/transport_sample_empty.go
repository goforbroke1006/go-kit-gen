package testdata

import (
	"github.com/go-kit/kit/transport/grpc"
)

type grpcXAXAXAXAServer struct {
	handleMethodOne grpc.Handler
	handleMethodTwo grpc.Handler
	handleGetOdds   grpc.Handler
}
