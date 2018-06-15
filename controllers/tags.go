package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//TagGet handles GET /tags/:title route
func TagGet(c *gin.Context) {
	db := models.GetDB()
	tag := models.Tag{}
	db.Preload("Posts", "published = true").Preload("Posts.Comments", "published = true").Preload("Posts.Tags").Preload("Posts.User").First(&tag, c.Param("title"))
	if len(tag.Title) == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = tag.Title
	h["Tag"] = tag
	c.HTML(http.StatusOK, "tags/show", h)
}

//TagIndex handles GET /admin/tags route
func TagIndex(c *gin.Context) {
	db := models.GetDB()
	var tags []models.Tag
	db.Preload("Posts").Order("title asc").Find(&tags)
	h := DefaultH(c)
	h["Title"] = "List of tags"
	h["Tags"] = tags
	c.HTML(http.StatusOK, "tags/index", h)
}

//TagNew handles GET /admin/new_tag route
func TagNew(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "New tag"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "tags/form", h)
}

//TagCreate handles POST /admin/new_tag route
func TagCreate(c *gin.Context) {
	tag := models.Tag{}
	db := models.GetDB()
	if err := c.ShouldBind(&tag); err != nil {
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

//TagDelete handles POST /admin/tags/:title/delete route
func TagDelete(c *gin.Context) {
	db := models.GetDB()
	tag := models.Tag{}
	db.First(&tag, c.Param("title"))
	if len(tag.Title) == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	if err := db.Delete(&tag).Error; err != nil {
		logrus.Error(err)
		c.HTML(http.StatusInternalServerError, "errors/500", gin.H{"Error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/admin/tags")
}
