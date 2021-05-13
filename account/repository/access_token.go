package repository

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

var createAccessTokenError = errors.New("Create AccessToken failed")
var createOrUpdateAccessTokenError = errors.New("Create or update AccessToken failed")
var verifyAccessTokenError = errors.New("Token verification failed")
var verifyAccessTokenNotValidError = errors.New("Token is not valid")
var verifyAccessTokenDoesNotMatchError = errors.New("Token doesn't match")

type UserAccessToken struct {
	gorm.Model
	User        User `gorm:"foreignkey:user_id;association_foreignkey:id"` // use UserRefer as foreign key
	UserID      uint
	AccessToken string `gorm:"type:varchar(255);unique_index;not null" json:"access_token"`
}

// for generating token
type claims struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

func generateAccessToken(db *gorm.DB, userID uint) (string, error) {
	claims := &claims{
		UserId: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour * ExpirationTimeDays).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(env.GetEnvItem("SECRET_KEY")))

	if err != nil {
		return "", createAccessTokenError
	}

	var tokenModel = UserAccessToken{}
	tokenModel.UserID = userID
	tokenModel.AccessToken = tokenString

	//UpdateOrCreate
	err = db.Scopes(scopeUser(tokenModel.UserID)).Assign(UserAccessToken{AccessToken: tokenModel.AccessToken}).FirstOrCreate(tokenModel).Error
	if err != nil {
		return "", createOrUpdateAccessTokenError
	}

	return tokenModel.AccessToken, nil

}

func verifyAccessToken(db *gorm.DB, token string) (uint, error) {
	claims := &claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return 0, verifyAccessTokenError
	}

	if !tkn.Valid {
		return 0, verifyAccessTokenNotValidError
	}
	var tokenModel = UserAccessToken{}

	if result := db.Scopes(scopeToken(token)).Find(&tokenModel); result.Error != nil {
		return 0, verifyAccessTokenDoesNotMatchError
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
