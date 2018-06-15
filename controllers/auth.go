package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//SignInGet handles GET /signin route
func SignInGet(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "Basic GIN web-site signin form"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "auth/signin", h)
}

//SignInPost handles POST /signin route, authenticates user
func SignInPost(c *gin.Context) {
	session := sessions.Default(c)
	login := models.Login{}
	db := models.GetDB()
	returnURL := c.DefaultQuery("return", "/admin/")
	if err := c.ShouldBind(&login); err != nil {
		session.AddFlash("Please, fill out form correctly.")
		session.Save()
		c.Redirect(http.StatusFound, "/signin")
		return
	}

	user := models.User{}
	db.Where("email = lower(?)", login.Email).First(&user)

	if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)) != nil {
		logrus.Errorf("Login error, IP: %s, Email: %s", c.ClientIP(), login.Email)
		session.AddFlash("Email or password incorrect")
		session.Save()
		c.Redirect(http.StatusFound, fmt.Sprintf("/signin?return=%s", url.QueryEscape(returnURL)))
		return
	}

	session.Set(userIDKey, user.ID)
	session.Save()
	c.Redirect(http.StatusFound, returnURL)
}

//SignUpGet handles GET /signup route
func SignUpGet(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "Basic GIN web-site signup form"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "auth/signup", h)
}

//SignUpPost handles POST /signup route, creates new user
func SignUpPost(c *gin.Context) {
	session := sessions.Default(c)
	register := models.Register{}
	db := models.GetDB()
	if err := c.ShouldBind(&register); err != nil {
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/signup")
		return
	}
	register.Email = strings.ToLower(register.Email)
	user := models.User{}
	db.Where("email = ?", register.Email).First(&user)
	if user.ID != 0 {
		session.AddFlash("User exists")
		session.Save()
		c.Redirect(http.StatusFound, "/signup")
		return
	}
	//create user
	user.Email = register.Email
	user.Password = register.Password
	user.Name = register.Name
	if err := db.Create(&user).Error; err != nil {
		session.AddFlash("Error whilst registering user.")
		session.Save()
		logrus.Errorf("Error whilst registering user: %v", err)
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	session.Set(userIDKey, user.ID)
	session.Save()
	c.Redirect(http.StatusFound, "/")
	return
}

//LogoutGet handles GET /logout route
func LogoutGet(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(userIDKey)
	session.Save()
	c.Redirect(http.StatusSeeOther, "/")
}
