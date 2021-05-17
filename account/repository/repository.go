package repository

import (
	"context"
	"github.com/amiranmanesh/go-smart-api-maker/account/service"
	"github.com/amiranmanesh/go-smart-api-maker/db/sql"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gorm.io/gorm"
)

func NewAccountRepository(db *gorm.DB, logger log.Logger) service.IRepository {
	if err := db.AutoMigrate(sql.User{}, sql.UserAccessToken{}); err != nil {
		level.Error(logger).Log("Repository auto migration failed", err)
		panic(err)
	}
	return &repository{db, log.With(logger, "Repository")}
}

type repository struct {
	db     *gorm.DB
	logger log.Logger
}

func (r repository) SignUp(ctx context.Context, name, email, password string) (string, error) {
	logger := log.With(r.logger, "SignUp")
	logger.Log("Start")

	user := &sql.User{}
	user.Name = name
	user.Email = email
	user.Password = password

	if err := user.Save(r.db); err != nil {
		level.Error(logger).Log("Error is: ", err)
		return "", err
	}

	token, err := sql.GenerateAccessToken(r.db, user.ID)
	if err != nil {
		level.Error(logger).Log("Error is: ", err)
		return "", err
	}

	return token, nil
}

func (r repository) Login(ctx context.Context, email, password string) (string, error) {
	logger := log.With(r.logger, "Login")
	logger.Log("Start")

	user := &sql.User{}
	user.Email = email
	user.Password = password

	if err := user.Login(r.db); err != nil {
		level.Error(logger).Log("Error is: ", err)
		return "", err
	}

	token, err := sql.GenerateAccessToken(r.db, user.ID)
	if err != nil {
		level.Error(logger).Log("Error is: ", err)
		return "", err
	}

	return token, nil
}

func (r repository) Verify(ctx context.Context, token string) (uint, error) {
	logger := log.With(r.logger, "Verify")
	logger.Log("Start")

	uid, err := sql.VerifyAccessToken(r.db, token)
	if err != nil {
		level.Error(logger).Log("Error is: ", err)
		return 0, err
	}

	model := &sql.User{}
	model.ID = uid
	if err := model.Find(r.db); err != nil {
		level.Error(logger).Log("Error is: ", err)
		return 0, err
	}

	return model.ID, nil
}
