package admin

import (
	"net/http"

	"github.com/denisbakhtin/ginblog/helpers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GET tag list
func TagIndex(c *gin.Context) {
	list, err := models.GetTags()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "List of tags"
	h["List"] = list
	h["Active"] = "tags"
	c.HTML(http.StatusOK, "admin/tags/index", h)
}

// GET tag creation form
func TagNew(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "New tag"
	h["Active"] = "tags"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "admin/tags/form", h)
}

// tag tag creation form
func TagCreate(c *gin.Context) {
	tag := &models.Tag{}
	if err := c.Bind(tag); err == nil {
		if err := tag.Insert(); err != nil {
			c.HTML(http.StatusInternalServerError, "errors/500", nil)
			return
		}
		c.Redirect(http.StatusFound, "/admin/tags")
	} else {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_tag")
	}
}

// tag tag deletion request
func TagDelete(c *gin.Context) {
	tag, _ := models.GetTag(c.Param("name"))
	if err := tag.Delete(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		return
	} else {
		c.Redirect(http.StatusFound, "/admin/tags")
	}
}
