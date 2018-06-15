package controllers

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//PostGet handles GET /posts/:id route
func PostGet(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)
	post := models.Post{}
	db.Preload("Tags").Preload("User").Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("comments.created_at DESC")
	}).First(&post, c.Param("id"))

	if post.ID == 0 || !post.Published {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = post.Title
	h["Post"] = post
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "posts/show", h)
}

//PostIndex handles GET /admin/posts route
func PostIndex(c *gin.Context) {
	db := models.GetDB()
	var posts []models.Post
	db.Preload("Tags").Find(&posts)
	h := DefaultH(c)
	h["Title"] = "List of blog posts"
	h["Posts"] = posts
	c.HTML(http.StatusOK, "posts/index", h)
}

//PostNew handles GET /admin/new_post route
func PostNew(c *gin.Context) {
	var tags []models.Tag
	db := models.GetDB()
	db.Order("title asc").Find(&tags)
	h := DefaultH(c)
	h["Title"] = "New post"
	h["Tags"] = tags
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "posts/form", h)
}

//PostCreate handles POST /admin/new_post route
func PostCreate(c *gin.Context) {
	post := models.Post{}
	db := models.GetDB()
	if err := c.ShouldBind(&post); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		logrus.Error(err)
		c.Redirect(http.StatusSeeOther, "/admin/new_post")
		return
	}
	tags := make([]models.Tag, 0, len(post.FormTags))
	for i := range post.FormTags {
		tags = append(tags, models.Tag{Title: post.FormTags[i]})
	}
	post.Tags = tags
	if user, exists := c.Get("User"); exists {
		post.UserID = user.(*models.User).ID
	}
	if err := db.Create(&post).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/posts")
}

//PostEdit handles GET /admin/posts/:id/edit route
func PostEdit(c *gin.Context) {
	db := models.GetDB()
	post := models.Post{}
	db.Preload("Tags").First(&post, c.Param("id"))
	if post.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = "Edit post entry"
	h["Post"] = post
	h["Tags"] = post.Tags
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "posts/form", h)
}

//PostUpdate handles POST /admin/posts/:id/edit route
func PostUpdate(c *gin.Context) {
	db := models.GetDB()
	post := models.Post{}
	if err := c.ShouldBind(&post); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		logrus.Error(err)
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/posts/%s/edit", c.Param("id")))
		return
	}
	tags := make([]models.Tag, 0, len(post.FormTags))
	for i := range post.FormTags {
		tags = append(tags, models.Tag{Title: post.FormTags[i]})
	}
	post.Tags = tags

	if err := db.Save(&post).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	if err := db.Exec("DELETE FROM posts_tags WHERE post_id = ? AND tag_title NOT IN(?)", post.ID, post.FormTags).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/posts")
}

//PostDelete handles POST /admin/posts/:id/delete route
func PostDelete(c *gin.Context) {
	db := models.GetDB()
	post := models.Post{}
	db.First(&post, c.Param("id"))
	if post.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	if err := db.Delete(&post).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/posts")
}
