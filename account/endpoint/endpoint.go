package endpoint

import (
	"context"
	"github.com/amiranmanesh/go-smart-api-maker/account/layers"
	"github.com/amiranmanesh/go-smart-api-maker/account/server"
	"github.com/go-kit/kit/endpoint"
)

func MakeEndpoint(s layers.Service) layers.Endpoints {
	return layers.Endpoints{
		SignUp: makeSignUpEndpoint(s),
		Login:  makeLoginEndpoint(s),
		Verify: makeVerifyEndpoint(s),
	}
}

func makeSignUpEndpoint(s layers.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(server.SignUpRequest)
		token, err := s.SignUp(ctx, req.Name, req.Email, req.Password)
		if err != nil {
			return server.SignUpResponse{
				Success: false,
				Token:   nil,
			}, err
		} else {
			return server.SignUpResponse{
				Success: true,
				Token:   token,
			}, nil
		}
	}
}

func makeLoginEndpoint(s layers.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(server.LoginRequest)
		token, err := s.Login(ctx, req.Email, req.Password)
		if err != nil {
			return server.LoginResponse{
				Success: false,
				Token:   nil,
			}, err
		} else {
			return server.LoginResponse{
				Success: true,
				Token:   token,
			}, nil
		}
	}
}

func makeVerifyEndpoint(s layers.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(server.VerifyRequest)
		userModel, err := s.Verify(ctx, req.Token)
		if err != nil {
			return server.VerifyResponse{
				Success: false,
				UserID:  nil,
			}, err
		} else {
			return server.VerifyResponse{
				Success: true,
				UserID:  userModel.ID,
			}, nil
		}
	}
}
