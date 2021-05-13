package logic

import (
	"context"
	"github.com/amiranmanesh/go-smart-api-maker/account/layers"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type repo struct {
	db     *gorm.DB
	logger log.Logger
}

func NewRepository(db *gorm.DB, logger log.Logger) layers.Repository {
	if err := db.AutoMigrate(User{}, AccessToken{}); err != nil {
		_ = logger.Log("repository auto migration failed", err)
		panic(err)
	}
	return &repo{db, log.With(logger, "repository", "db")}
}

func (r repo) SignUp(ctx context.Context, user User) (string, error) {
	_ = r.logger.Log("SignUp")

	if err := user.Save(r.db); err != nil {
		_ = r.logger.Log("SignUp", "Error is: ", err)
		return "", err
	}

	token, err := generateAccessToken(r.db, user.ID)
	if err != nil {
		_ = r.logger.Log("SignUp", "Error is: ", err)
		return "", err
	}

	return token, nil
}

func (r repo) Login(ctx context.Context, user User) (string, error) {
	_ = r.logger.Log("Login")

	if err := user.Login(r.db); err != nil {
		_ = r.logger.Log("Login", "Error is: ", err)
		return "", err
	}

	token, err := generateAccessToken(r.db, user.ID)
	if err != nil {
		_ = r.logger.Log("Login", "Error is: ", err)
		return "", err
	}

	return token, nil
}

func (r repo) Verify(ctx context.Context, token string) (*User, error) {
	_ = r.logger.Log("Verify")

	uid, err := verifyAccessToken(r.db, token)
	if err != nil {
		_ = r.logger.Log("Verify", "Error is: ", err)
		return nil, err
	}

	model := &User{}
	model.ID = uid
	if err := model.Find(r.db); err != nil {
		_ = r.logger.Log("Verify", "Error is: ", err)
		return nil, err
	}

	return model, nil
}
