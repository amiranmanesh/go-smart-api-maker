package server

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type IService interface {
	SignUp(ctx context.Context, name, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	Verify(ctx context.Context, token string) (uint, error)
}

func MakeEndpoint(s IService) Endpoints {
	return Endpoints{
		SignUp:      makeSignUpEndpoint(s),
		Login:       makeLoginEndpoint(s),
		VerifyToken: makeVerifyTokenEndpoint(s),
	}
}

type Endpoints struct {
	SignUp      endpoint.Endpoint
	Login       endpoint.Endpoint
	VerifyToken endpoint.Endpoint
}

func makeSignUpEndpoint(s IService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignUpRequest)
		token, err := s.SignUp(ctx, req.Name, req.Email, req.Password)
		if err != nil {
			return SignUpResponse{Success: false}, err
		} else {
			return SignUpResponse{
				Success: true,
				Token:   token,
			}, nil
		}
	}
}

func makeLoginEndpoint(s IService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		token, err := s.Login(ctx, req.Email, req.Password)
		if err != nil {
			return LoginResponse{Success: false}, err
		} else {
			return LoginResponse{
				Success: true,
				Token:   token,
			}, nil
		}
	}
}

func makeVerifyTokenEndpoint(s IService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyTokenRequest)
		userID, err := s.Verify(ctx, req.Token)
		if err != nil {
			return VerifyTokenResponse{Success: false}, err
		} else {
			return VerifyTokenResponse{
				Success: true,
				UserID:  int32(userID),
			}, nil
		}
	}
}
