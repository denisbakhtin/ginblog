package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/helpers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"
)

//TagGet handles GET /tags/:name route
func TagGet(c *gin.Context) {
	db := models.GetDB()
	tag := models.Tag{}
	db.Where("lower(name) = $1", strings.ToLower(c.Param("name"))).First(&tag)
	if tag.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	var list []models.Post
	db.Where("published = true AND tag_id = $1", tag.ID).Find(&list)
	h := helpers.DefaultH(c)
	h["Title"] = tag.Name
	h["Active"] = "tags"
	h["List"] = list
	c.HTML(http.StatusOK, "tags/show", h)
}

//TagIndex handles GET /admin/tags route
func TagIndex(c *gin.Context) {
	db := models.GetDB()
	var list []models.Tag
	db.Find(&list)
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
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "tags/form", h)
}

//TagCreate handles POST /admin/new_tag route
func TagCreate(c *gin.Context) {
	tag := &models.Tag{}
	db := models.GetDB()
	if err := c.Bind(tag); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_tag")
		return
	}

	if err := db.Create(&tag).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/tags")
}

//TagDelete handles POST /admin/tags/:name/delete route
func TagDelete(c *gin.Context) {
	db := models.GetDB()
	tag := models.Tag{}
	db.Where("lower(name) = $1", strings.ToLower(c.Param("name"))).First(&tag)
	if tag.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	if err := db.Delete(&tag); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/tags")
}
