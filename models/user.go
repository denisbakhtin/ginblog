package models

import (
	"golang.org/x/crypto/bcrypt"
)

//Login view model
type Login struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

//Register view model
type Register struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

//User type contains user info
type User struct {
	Model

	Email    string `form:"email" binding:"required"`
	Name     string `form:"name"`
	Password string `form:"password" binding:"required"`
}

//BeforeSave gorm hook
func (u *User) BeforeSave() (err error) {
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	u.Password = string(hash)
	return
}
