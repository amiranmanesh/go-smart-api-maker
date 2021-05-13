package layers

import (
	"context"
	"github.com/amiranmanesh/go-smart-api-maker/account/logic"
)

// Repository describes the persistence on user model
type Repository interface {
	SignUp(ctx context.Context, user logic.User) (string, error)
	Login(ctx context.Context, user logic.User) (string, error)
	Verify(ctx context.Context, token string) (*logic.User, error)
}
