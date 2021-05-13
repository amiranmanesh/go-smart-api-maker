package service

import (
	"context"
	"github.com/amiranmanesh/go-smart-api-maker/account/layers"
	"github.com/amiranmanesh/go-smart-api-maker/account/logic"
	"github.com/amiranmanesh/go-smart-api-maker/utils/encrypting"
	"github.com/go-kit/kit/log"
)

type service struct {
	repository layers.Repository
	logger     log.Logger
}

func NewService(repository layers.Repository, logger log.Logger) layers.Service {
	return &service{repository, log.With(logger, "service")}
}

func (s service) SignUp(ctx context.Context, name, email, password string) (string, error) {

	model := logic.User{}
	model.Name = name
	model.Email = email
	model.Password = encrypting.GetHashedPassword(password)

	return s.repository.SignUp(ctx, model)
}

func (s service) Login(ctx context.Context, email, password string) (string, error) {

	model := logic.User{}
	model.Email = email
	model.Password = encrypting.GetHashedPassword(password)

	return s.repository.Login(ctx, model)
}

func (s service) Verify(ctx context.Context, token string) (*logic.User, error) {
	return s.Verify(ctx, token)
}
