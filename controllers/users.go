package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//UserIndex handles GET /admin/users route
func UserIndex(c *gin.Context) {
	db := models.GetDB()
	var users []models.User
	db.Find(&users)
	h := DefaultH(c)
	h["Title"] = "List of users"
	h["Users"] = users
	c.HTML(http.StatusOK, "users/index", h)
}

//UserNew handles GET /admin/new_user route
func UserNew(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "New user"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "users/form", h)
}

//UserCreate handles POST /admin/new_user route
func UserCreate(c *gin.Context) {
	user := &models.User{}
	db := models.GetDB()
	if err := c.ShouldBind(user); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_user")
		return
	}

	if err := db.Create(&user).Error; err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_user")
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}

//UserEdit handles GET /admin/users/:id/edit route
func UserEdit(c *gin.Context) {
	db := models.GetDB()
	user := models.User{}
	db.First(&user, c.Param("id"))
	if user.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = "Edit user"
	h["User"] = user
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "users/form", h)
}

//UserUpdate handles POST /admin/users/:id/edit route
func UserUpdate(c *gin.Context) {
	user := &models.User{}
	db := models.GetDB()
	if err := c.ShouldBind(user); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/users")
		return
	}

	if err := db.Save(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}

//UserDelete handles POST /admin/users/:id/delete route
func UserDelete(c *gin.Context) {
	db := models.GetDB()
	user := models.User{}
	db.First(&user, c.Param("id"))
	if user.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	if err := db.Delete(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}
