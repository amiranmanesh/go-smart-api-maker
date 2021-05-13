package layers

import (
	"context"
	"github.com/amiranmanesh/go-smart-api-maker/account/logic"
	"github.com/go-kit/kit/endpoint"
)

//Endpoint Layer
type Endpoints struct {
	SignUp endpoint.Endpoint
	Login  endpoint.Endpoint
	Verify endpoint.Endpoint
}

//Logic Layer
type Service interface {
	SignUp(ctx context.Context, name, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	Verify(ctx context.Context, token string) (*logic.User, error)
}

type Repository interface {
	SignUp(ctx context.Context, user logic.User) (string, error)
	Login(ctx context.Context, user logic.User) (string, error)
	Verify(ctx context.Context, token string) (*logic.User, error)
}
