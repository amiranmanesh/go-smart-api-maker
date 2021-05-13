package logic

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

var signUpEmailExistError = errors.New("Email is already exist")
var signUpCreateUserError = errors.New("Create User failed")
var loginUserNotFoundError = errors.New("User not found")
var loginPasswordError = errors.New("Password is not valid")
var userNotFoundError = errors.New("User not found")

type User struct {
	gorm.Model
	Name            string    `gorm:"type:varchar(100); not null" json:"name"`
	Email           string    `gorm:"type:varchar(100);unique_index; not null" json:"email"`
	Password        string    `json:"-"`
	EmailVerifiedAt time.Time `gorm:"type:varchar(100)" json:"-"`
}

func (u *User) Save(db *gorm.DB) error {

	result := db.Scopes(scopeEmail(u.Email)).First(&u)
	if err := result.Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return signUpEmailExistError
	}

	if err := db.Create(&u).Error; err != nil {
		return signUpCreateUserError
	}

	return nil

}

func (u *User) Login(db *gorm.DB) error {

	inputPass := u.Password

	isFounded := db.Scopes(scopeEmail(u.Email)).Find(&u)
	if isFounded.Error != nil {
		return loginUserNotFoundError
	}

	if u.Password != inputPass {
		return loginPasswordError
	}

	return nil
}
func (u *User) Find(db *gorm.DB) error {
	result := db.First(&u, u.ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return userNotFoundError
	}
	return nil
}

//scopes
func scopeEmail(email string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	}
}
