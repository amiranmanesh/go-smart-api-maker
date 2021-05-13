package server

import (
	"context"
	"github.com/amiranmanesh/go-smart-api-maker/account/service"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	SignUp endpoint.Endpoint
	Login  endpoint.Endpoint
	Verify endpoint.Endpoint
}

func MakeEndpoint(s service.Service) Endpoints {
	return Endpoints{
		SignUp: makeSignUpEndpoint(s),
		Login:  makeLoginEndpoint(s),
		Verify: makeVerifyEndpoint(s),
	}
}

func makeSignUpEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignUpRequest)
		token, err := s.SignUp(ctx, req.Name, req.Email, req.Password)
		if err != nil {
			return SignUpResponse{
				Success: false,
				Token:   nil,
			}, err
		} else {
			return SignUpResponse{
				Success: true,
				Token:   &token,
			}, nil
		}
	}
}

func makeLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		token, err := s.Login(ctx, req.Email, req.Password)
		if err != nil {
			return LoginResponse{
				Success: false,
				Token:   nil,
			}, err
		} else {
			return LoginResponse{
				Success: true,
				Token:   &token,
			}, nil
		}
	}
}

func makeVerifyEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyRequest)
		userModel, err := s.Verify(ctx, req.Token)
		if err != nil {
			return VerifyResponse{
				Success: false,
				UserID:  nil,
			}, err
		} else {
			return VerifyResponse{
				Success: true,
				UserID:  &userModel.ID,
			}, nil
		}
	}
}
