package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginbasic/helpers"
	"github.com/denisbakhtin/ginbasic/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GET /signin route
func SignInGet(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "Basic GIN web-site signin form"
	h["Active"] = "signin"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "auth/signin", h)
}

// POST /signin route, authenticates user
func SignInPost(c *gin.Context) {
	session := sessions.Default(c)
	user := &models.User{}
	if err := c.Bind(user); err == nil {
		userDB, _ := models.GetUserByEmail(user.Email)
		if userDB.Id == 0 {
			logrus.Errorf("Login error, IP: %s, Email: %s", c.ClientIP(), user.Email)
			session.AddFlash("Email or password incorrect")
			session.Save()
			c.Redirect(http.StatusFound, "/signin")
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password)); err != nil {
			logrus.Errorf("Login error, IP: %s, Email: %s", c.ClientIP(), user.Email)
			session.AddFlash("Email or password incorrect")
			session.Save()
			c.Redirect(http.StatusFound, "/signin")
			return
		}

		session.Set("UserId", userDB.Id)
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return

	} else {
		session.AddFlash("Please, fill out form correctly.")
		session.Save()
		c.Redirect(http.StatusFound, "/signin")
		return
	}
}

// GET /signup route
func SignUpGet(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "Basic GIN web-site signup form"
	h["Active"] = "signup"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "auth/signup", h)
}

// POST /signup route, creates new user
func SignUpPost(c *gin.Context) {
	session := sessions.Default(c)
	user := &models.User{}
	if err := c.Bind(user); err == nil {
		userDB, _ := models.GetUserByEmail(user.Email)
		if userDB.Id != 0 {
			session.AddFlash("User exists")
			session.Save()
			c.Redirect(http.StatusFound, "/signup")
			return
		}
		//create user
		err := user.HashPassword()
		if err != nil {
			session.AddFlash("Error whilst registering user.")
			session.Save()
			logrus.Errorf("Error whilst registering user: %v", err)
			c.Redirect(http.StatusFound, "/signup")
			return
		}

		if err := user.Insert(); err != nil {
			session.AddFlash("Error whilst registering user.")
			session.Save()
			logrus.Errorf("Error whilst registering user: %v", err)
			c.Redirect(http.StatusFound, "/signup")
			return
		}

		session.Set("UserId", user.Id)
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return

	} else {
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/signup")
		return
	}
}

// GET /logout route
func LogoutGet(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("UserId")
	session.Save()
	c.Redirect(http.StatusSeeOther, "/")
}
