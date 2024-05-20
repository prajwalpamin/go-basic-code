package model

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model `json:"-,omitempty"`
	Name       string `json:"name" gorm:"unique"`
	Email      string `json:"email" gorm:"unique"`
	Password   string `json:"password"`
}
type LoginUser struct {
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

//reciever user
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) ComparePassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
