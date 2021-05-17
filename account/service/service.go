package service

import (
	"context"
	"github.com/amiranmanesh/go-smart-api-maker/account/server"
	"github.com/amiranmanesh/go-smart-api-maker/utils/encrypting"
	"github.com/go-kit/kit/log"
)

type IRepository interface {
	SignUp(ctx context.Context, name, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	Verify(ctx context.Context, token string) (uint, error)
}

func NewService(repository IRepository, logger log.Logger) server.IService {
	return &service{repository, log.With(logger, "service")}
}

type service struct {
	repository IRepository
	logger     log.Logger
}

func (s service) SignUp(ctx context.Context, name, email, password string) (string, error) {
	hashPassword := encrypting.GetHashedPassword(password)
	return s.repository.SignUp(ctx, name, email, hashPassword)
}

func (s service) Login(ctx context.Context, email, password string) (string, error) {
	hashPassword := encrypting.GetHashedPassword(password)
	return s.repository.Login(ctx, email, hashPassword)
}

func (s service) Verify(ctx context.Context, token string) (uint, error) {
	return s.repository.Verify(ctx, token)
}
