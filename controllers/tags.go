package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/helpers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/gin"
)

//TagGet handles GET /tags/:name route
func TagGet(c *gin.Context) {
	tag, err := models.GetTag(c.Param("name"))
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	list, err := models.GetPostsByTag(tag.Name)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}

	h := helpers.DefaultH(c)
	h["Title"] = tag.Name
	h["Active"] = "tags"
	h["List"] = list
	c.HTML(http.StatusOK, "tags/show", h)
}

//TagIndex handles GET /admin/tags route
func TagIndex(c *gin.Context) {
	list, err := models.GetTags()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "List of tags"
	h["List"] = list
	h["Active"] = "tags"
	c.HTML(http.StatusOK, "tags/index", h)
}

//TagNew handles GET /admin/new_tag route
func TagNew(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "New tag"
	h["Active"] = "tags"
	session := helpers.GetSession(c)
	h["Flash"] = session.Flashes()
	session.Save(c.Request, c.Writer)

	c.HTML(http.StatusOK, "tags/form", h)
}

//TagCreate handles POST /admin/new_tag route
func TagCreate(c *gin.Context) {
	tag := &models.Tag{}
	if err := c.Bind(tag); err != nil {
		session := helpers.GetSession(c)
		session.AddFlash(err.Error())
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusSeeOther, "/admin/new_tag")
		return
	}

	if err := tag.Insert(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/tags")
}

//TagDelete handles POST /admin/tags/:name/delete route
func TagDelete(c *gin.Context) {
	tag, _ := models.GetTag(c.Param("name"))
	if err := tag.Delete(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/tags")
}
