package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/helpers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/gin"
)

//UserIndex handles GET /admin/users route
func UserIndex(c *gin.Context) {
	list, err := models.GetUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "List of users"
	h["Active"] = "users"
	h["List"] = list
	c.HTML(http.StatusOK, "users/index", h)
}

//UserNew handles GET /admin/new_user route
func UserNew(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "New user"
	h["Active"] = "users"
	session := helpers.GetSession(c)
	h["Flash"] = session.Flashes()
	session.Save(c.Request, c.Writer)
	c.HTML(http.StatusOK, "users/form", h)
}

//UserCreate handles POST /admin/new_user route
func UserCreate(c *gin.Context) {
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		session := helpers.GetSession(c)
		session.AddFlash(err.Error())
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusSeeOther, "/admin/new_user")
		return
	}

	if err := user.HashPassword(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	if err := user.Insert(); err != nil {
		session := helpers.GetSession(c)
		session.AddFlash(err.Error())
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusSeeOther, "/admin/new_user")
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}

//UserEdit handles GET /admin/users/:id/edit route
func UserEdit(c *gin.Context) {
	user, _ := models.GetUser(c.Param("id"))
	if user.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "Edit user"
	h["Active"] = "users"
	h["User"] = user
	session := helpers.GetSession(c)
	h["Flash"] = session.Flashes()
	session.Save(c.Request, c.Writer)
	c.HTML(http.StatusOK, "users/form", h)
}

//UserUpdate handles POST /admin/users/:id/edit route
func UserUpdate(c *gin.Context) {
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		session := helpers.GetSession(c)
		session.AddFlash(err.Error())
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusSeeOther, "/admin/users")
		return
	}

	if err := user.HashPassword(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	if err := user.Update(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}

//UserDelete handles POST /admin/users/:id/delete route
func UserDelete(c *gin.Context) {
	user, _ := models.GetUser(c.Param("id"))
	if err := user.Delete(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}
