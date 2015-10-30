package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64     `form:"id" json:"id"`
	Email     string    `form:"email" json:"email" binding:"required"`
	Name      string    `form:"name" json:"name"`
	Password  string    `form:"password" json:"password" binding:"required"`
	Timestamp time.Time `form:"-" json:"timestamp"`
}

func (user *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}

func (user *User) Insert() error {
	err := db.QueryRow("INSERT INTO users(email, name, password, timestamp) VALUES(lower($1),$2,$3,$4) RETURNING id", user.Email, user.Name, user.Password, time.Now()).Scan(&user.Id)
	return err
}

func (user *User) Update() error {
	_, err := db.Exec("UPDATE users SET email=lower($2), name=$3, password=$4 WHERE id=$1", user.Id, user.Email, user.Name, user.Password)
	return err
}

func (user *User) Delete() error {
	count := 0
	_ = db.Get(&count, "SELECT count(id) FROM users")
	if count <= 1 {
		return fmt.Errorf("Can't remove last user")
	}
	_, err := db.Exec("DELETE FROM users WHERE id=$1", user.Id)
	return err
}

func GetUser(id interface{}) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}

func GetUsers() ([]User, error) {
	var list []User
	err := db.Select(&list, "SELECT * FROM users ORDER BY id")
	return list, err
}

func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE lower(email)=lower($1)", email)
	return user, err
}
