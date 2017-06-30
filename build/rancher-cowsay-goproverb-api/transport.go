package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func makeTextsayEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(serviceRequest)
		index, say := svc.Textsay(req.I)
		return serviceResponse{index, say}, nil
	}
}

func makeCowsayEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(serviceRequest)
		index, say := svc.Cowsay(req.I)
		return serviceResponse{index, say}, nil
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request serviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return serviceRequest{}, nil
	}
	return request, nil
}

func encodeGenericJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeCowsayTextResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(serviceResponse)
	fmt.Fprint(w, res.S)
	return nil
}

type serviceRequest struct {
	I int `json:"index,omitempty"`
}

type serviceResponse struct {
	I int    `json:"index"`
	S string `json:"say"`
}
