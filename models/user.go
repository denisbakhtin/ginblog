package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//User type contains user info
type User struct {
	ID        int64     `form:"id" json:"id" database:"id"`
	Email     string    `form:"email" json:"email"`
	Name      string    `form:"name" json:"name"`
	Password  string    `form:"password" json:"password"`
	Timestamp time.Time `form:"-" json:"timestamp"`
}

//HashPassword substitutes User.Password with its bcrypt hash
func (user *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}

//Insert saves user info into db
func (user *User) Insert() error {
	err := db.QueryRow("INSERT INTO users(email, name, password, timestamp) VALUES(lower($1),$2,$3,$4) RETURNING id", user.Email, user.Name, user.Password, time.Now()).Scan(&user.ID)
	return err
}

//Update updates user info in db
func (user *User) Update() error {
	_, err := db.Exec("UPDATE users SET email=lower($2), name=$3, password=$4 WHERE id=$1", user.ID, user.Email, user.Name, user.Password)
	return err
}

//Delete removes user record from db. Can't remove the last user
func (user *User) Delete() error {
	count := 0
	_ = db.Get(&count, "SELECT count(id) FROM users")
	if count <= 1 {
		return fmt.Errorf("Can't remove last user")
	}
	_, err := db.Exec("DELETE FROM users WHERE id=$1", user.ID)
	return err
}

//GetUser returns user by his id
func GetUser(id interface{}) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}

//GetUsers returns a list of user ordered by id
func GetUsers() ([]User, error) {
	var list []User
	err := db.Select(&list, "SELECT * FROM users ORDER BY id")
	return list, err
}

//GetUserByEmail returns user record by his email, case insensitive
func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE lower(email)=lower($1)", email)
	return user, err
}
