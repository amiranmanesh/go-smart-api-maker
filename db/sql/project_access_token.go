package sql

import (
	"errors"
	"github.com/amiranmanesh/go-smart-api-maker/utils/env"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"os"
)

var createProjectAccessTokenError = errors.New("create project access token failed")
var createOrUpdateProjectAccessTokenError = errors.New("create or update project access token failed")
var verifyProjectAccessTokenError = errors.New("project token verification failed")
var verifyProjectAccessTokenNotValidError = errors.New("project token is not valid")
var verifyProjectAccessTokenDoesNotMatchError = errors.New("project token doesn't match")

type ProjectAccessToken struct {
	gorm.Model
	Project      Project `gorm:"foreignkey:project_id;association_foreignkey:id"` // use UserRefer as foreign key
	ProjectID    uint
	ProjectToken string `gorm:"type:varchar(255);unique_index;not null" json:"project_access_token"`
}

// for generating token
type projectClaims struct {
	ProjectId uint `json:"project_id"`
	jwt.StandardClaims
}

func GenerateProjectAccessToken(db *gorm.DB, projectID uint) (string, error) {
	claims := &projectClaims{
		ProjectId:      projectID,
		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(env.GetEnvItem("SECRET_KEY")))

	if err != nil {
		return "", createProjectAccessTokenError
	}

	var tokenModel = ProjectAccessToken{}
	tokenModel.ProjectID = projectID
	tokenModel.ProjectToken = tokenString

	//UpdateOrCreate
	err = db.Scopes(scopeProjectID(tokenModel.ProjectID)).Assign(ProjectAccessToken{ProjectToken: tokenModel.ProjectToken}).FirstOrCreate(tokenModel).Error
	if err != nil {
		return "", createOrUpdateProjectAccessTokenError
	}

	return tokenModel.ProjectToken, nil

}

func VerifyProjectAccessToken(db *gorm.DB, token string) (uint, error) {
	claims := &userClaims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return 0, verifyProjectAccessTokenError
	}

	if !tkn.Valid {
		return 0, verifyProjectAccessTokenNotValidError
	}
	var tokenModel = UserAccessToken{}

	if result := db.Scopes(scopeProjectToken(token)).Find(&tokenModel); result.Error != nil {
		return 0, verifyProjectAccessTokenDoesNotMatchError
	}

	return tokenModel.UserID, nil
}

//scopes
func scopeProjectID(userId uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("project_id = ?", userId)
	}
}

func scopeProjectToken(token string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("project_access_token = ?", token)
	}
}
