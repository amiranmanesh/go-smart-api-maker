package logic

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gorm.io/gorm"
)

type Repository interface {
	SignUp(ctx context.Context, user User) (string, error)
	Login(ctx context.Context, user User) (string, error)
	Verify(ctx context.Context, token string) (*User, error)
}

func NewRepository(db *gorm.DB, logger log.Logger) Repository {
	if err := db.AutoMigrate(User{}, UserAccessToken{}); err != nil {
		level.Error(logger).Log("Repository auto migration failed", err)
		panic(err)
	}
	return &repo{db, log.With(logger, "Repository")}
}

type repo struct {
	db     *gorm.DB
	logger log.Logger
}

func (r repo) SignUp(ctx context.Context, user User) (string, error) {
	logger := log.With(r.logger, "SignUp")
	logger.Log("Start")

	if err := user.Save(r.db); err != nil {
		level.Error(logger).Log("Error is: ", err)
		return "", err
	}

	token, err := generateAccessToken(r.db, user.ID)
	if err != nil {
		level.Error(logger).Log("Error is: ", err)
		return "", err
	}

	return token, nil
}

func (r repo) Login(ctx context.Context, user User) (string, error) {
	logger := log.With(r.logger, "Login")
	logger.Log("Start")

	if err := user.Login(r.db); err != nil {
		level.Error(logger).Log("Error is: ", err)
		return "", err
	}

	token, err := generateAccessToken(r.db, user.ID)
	if err != nil {
		level.Error(logger).Log("Error is: ", err)
		return "", err
	}

	return token, nil
}

func (r repo) Verify(ctx context.Context, token string) (*User, error) {
	logger := log.With(r.logger, "Verify")
	logger.Log("Start")

	uid, err := verifyAccessToken(r.db, token)
	if err != nil {
		level.Error(logger).Log("Error is: ", err)
		return nil, err
	}

	model := &User{}
	model.ID = uid
	if err := model.Find(r.db); err != nil {
		level.Error(logger).Log("Error is: ", err)
		return nil, err
	}

	return model, nil
}
