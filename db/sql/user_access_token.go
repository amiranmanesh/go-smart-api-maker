package sql

import (
	"errors"
	"github.com/amiranmanesh/go-smart-api-maker/utils/env"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"os"
	"time"
)

const (
	ExpirationTimeDays = 90
)

var createUserAccessTokenError = errors.New("create user access token failed")
var createOrUpdateUserAccessTokenError = errors.New("create or update user access token failed")
var verifyUserAccessTokenError = errors.New("user token verification failed")
var verifyUserAccessTokenNotValidError = errors.New("user token is not valid")
var verifyUserAccessTokenDoesNotMatchError = errors.New("user token doesn't match")

type UserAccessToken struct {
	gorm.Model
	User        User `gorm:"foreignkey:user_id;association_foreignkey:id"` // use UserRefer as foreign key
	UserID      uint
	AccessToken string `gorm:"type:varchar(255);unique_index;not null" json:"access_token"`
}

// for generating token
type userClaims struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateUserAccessToken(db *gorm.DB, userID uint) (string, error) {
	claims := &userClaims{
		UserId: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour * ExpirationTimeDays).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(env.GetEnvItem("SECRET_KEY")))

	if err != nil {
		return "", createUserAccessTokenError
	}

	var tokenModel = UserAccessToken{}
	tokenModel.UserID = userID
	tokenModel.AccessToken = tokenString

	//UpdateOrCreate
	err = db.Scopes(scopeUser(tokenModel.UserID)).Assign(UserAccessToken{AccessToken: tokenModel.AccessToken}).FirstOrCreate(tokenModel).Error
	if err != nil {
		return "", createOrUpdateUserAccessTokenError
	}

	return tokenModel.AccessToken, nil

}

func VerifyUserAccessToken(db *gorm.DB, token string) (uint, error) {
	claims := &userClaims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return 0, verifyUserAccessTokenError
	}

	if !tkn.Valid {
		return 0, verifyUserAccessTokenNotValidError
	}
	var tokenModel = UserAccessToken{}

	if result := db.Scopes(scopeToken(token)).Find(&tokenModel); result.Error != nil {
		return 0, verifyUserAccessTokenDoesNotMatchError
	}

	return tokenModel.UserID, nil
}

//scopes
func scopeUser(userId uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userId)
	}
}

func scopeToken(token string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("access_token = ?", token)
	}
}
