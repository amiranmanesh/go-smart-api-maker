package sql

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

var signUpEmailExistError = errors.New("email is already exist")
var signUpCreateUserError = errors.New("create User failed")
var loginUserNotFoundError = errors.New("user not found")
var loginPasswordError = errors.New("password is not valid")
var userNotFoundError = errors.New("user not found")

type User struct {
	gorm.Model
	Name            string    `gorm:"type:varchar(100); not null" json:"name"`
	Email           string    `gorm:"type:varchar(100);primaryKey;unique_index;not null;" json:"email"`
	Password        string    `json:"-"`
	EmailVerifiedAt time.Time `gorm:"type:varchar(100)" json:"-"`
}

func (u *User) Save(db *gorm.DB) error {

	if result := db.First(&u, "email = ?", u.Email); result.Error != nil {
		return signUpEmailExistError
	}

	if result := db.Create(&u); result.Error != nil {
		return signUpCreateUserError
	}

	return nil
}

func (u *User) Login(db *gorm.DB) error {
	inputPass := u.Password

	if result := db.First(&u, "email = ?", u.Email); result.Error != nil {
		return loginUserNotFoundError
	}

	if u.Password != inputPass {
		return loginPasswordError
	}

	return nil
}

func (u *User) Find(db *gorm.DB) error {
	if result := db.First(&u, "id = ?", u.ID); result.Error != nil {
		return userNotFoundError
	}
	return nil
}
