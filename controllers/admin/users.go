package admin

import (
	"net/http"

	"github.com/denisbakhtin/ginbasic/helpers"
	"github.com/denisbakhtin/ginbasic/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GET user list
func UserIndex(c *gin.Context) {
	list, err := models.GetUserList()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "List of users"
	h["Active"] = "users"
	h["List"] = list
	c.HTML(http.StatusOK, "admin/users/index", h)
}

// GET user creation form
func UserNew(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "New user"
	h["Active"] = "users"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "admin/users/form", h)
}

// POST user creation form
func UserCreate(c *gin.Context) {
	user := &models.User{}
	if err := c.Bind(user); err == nil {
		if err := user.HashPassword(); err != nil {
			c.HTML(http.StatusInternalServerError, "errors/500", nil)
			return
		}
		if err := user.Insert(); err != nil {
			session := sessions.Default(c)
			session.AddFlash(err.Error())
			session.Save()
			c.Redirect(http.StatusSeeOther, "/admin/new_user")
			return
		}
		c.Redirect(http.StatusFound, "/admin/users")
	} else {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_user")
	}
}

// GET user update form
func UserEdit(c *gin.Context) {
	user, _ := models.GetUser(c.Param("id"))
	if user.Id == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "Edit user"
	h["Active"] = "users"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "admin/users/form", h)
}

// POST user update form
func UserUpdate(c *gin.Context) {
	user := &models.User{}
	if err := c.Bind(user); err == nil {
		if err := user.HashPassword(); err != nil {
			c.HTML(http.StatusInternalServerError, "errors/500", nil)
			return
		}
		if err := user.Update(); err != nil {
			c.HTML(http.StatusInternalServerError, "errors/500", nil)
			return
		}
		c.Redirect(http.StatusFound, "/admin/users")
	} else {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/users")
	}
}

// POST user deletion request
func UserDelete(c *gin.Context) {
	user, _ := models.GetUser(c.Param("id"))
	if err := user.Delete(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		return
	} else {
		c.Redirect(http.StatusFound, "/admin/users")
	}
}
