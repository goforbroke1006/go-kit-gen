package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type MethodOneRequest struct {
	// TODO: create fields
}

type MethodOneResponse struct {
	// TODO: create fields
	Err string
}

type SomeAwesomeHubEndpoints struct {
	MethodOneEndpoint endpoint.Endpoint
}

func MakeSomeAwesomeHubEndpoints(svc interface{}) SomeAwesomeHubEndpoints {
	return SomeAwesomeHubEndpoints{
		MethodOneEndpoint: makeMethodOneEndpoint(svc),
	}
}

func makeMethodOneEndpoint(svc interface{}) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//req := request.(ReceiveUpdatesRequest)
		//result, err := svc.ReceiveUpdates(ctx)
		//return ReceiveUpdatesResponse{field: result.field}, nil

		// TODO: write service call with request and transformation result to response

		return nil, nil
	}
}
