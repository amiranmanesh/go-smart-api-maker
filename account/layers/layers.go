package layers

import (
	"context"
	"github.com/amiranmanesh/go-smart-api-maker/account/logic"
)

type Service interface {
	SignUp(ctx context.Context, name, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	Verify(ctx context.Context, token string) (*logic.User, error)
}

// Repository describes the persistence on user model
type Repository interface {
	SignUp(ctx context.Context, user logic.User) (string, error)
	Login(ctx context.Context, user logic.User) (string, error)
	Verify(ctx context.Context, token string) (*logic.User, error)
}
