package testdata

import (
	"github.com/go-kit/kit/transport/grpc"
)

type grpcSomeAwesomeHubServer struct {
	handleMethodOne  grpc.Handler
	handleMethodTwo  grpc.Handler
	handleGetMatches grpc.Handler
	handleGetOdds    grpc.Handler
}
