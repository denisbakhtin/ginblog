package controllers

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//CommentGet handles GET /comments/:id route
func CommentGet(c *gin.Context) {
	db := models.GetDB()
	comment := models.Comment{}
	db.First(&comment, c.Param("id"))
	if comment.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = fmt.Sprintf("Comment #%d", comment.ID)
	h["Comment"] = comment
	c.HTML(http.StatusOK, "comments/show", h)
}

//CommentIndex handles GET /admin/comments route
func CommentIndex(c *gin.Context) {
	db := models.GetDB()
	var comments []models.Comment
	db.Find(&comments)
	h := DefaultH(c)
	h["Title"] = "List of comments"
	h["Comments"] = comments
	c.HTML(http.StatusOK, "comments/index", h)
}

//CommentNew handles GET /admin/new_comment route
func CommentNew(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "New comment"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "comments/form", h)
}

//CommentCreate handles POST /admin/new_comment route
func CommentCreate(c *gin.Context) {
	comment := models.Comment{}
	db := models.GetDB()
	if err := c.ShouldBind(&comment); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_comment")
		return
	}

	u, _ := c.Get("User")
	if user, ok := u.(*models.User); ok {
		comment.UserID = user.ID
		comment.UserName = user.Name
	}

	if err := db.Create(&comment).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/comments")
}

//CommentPublicCreate handles POST /new_comment route
func CommentPublicCreate(c *gin.Context) {
	comment := models.Comment{}
	db := models.GetDB()
	session := sessions.Default(c)

	if err := c.ShouldBind(&comment); err != nil {
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, c.Request.Referer())
		return
	}

	if name, ok := session.Get("oauth-username").(string); ok {
		comment.UserName = name
	}

	if err := db.Create(&comment).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	session.AddFlash("Your comment has been published!")
	session.Save()
	c.Redirect(http.StatusFound, c.Request.Referer())
}

//CommentEdit handles GET /admin/comments/:id/edit route
func CommentEdit(c *gin.Context) {
	db := models.GetDB()
	comment := models.Comment{}
	db.First(&comment, c.Param("id"))
	if comment.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = "Edit comment"
	h["Comment"] = comment
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "comments/form", h)
}

//CommentUpdate handles POST /admin/comments/:id/edit route
func CommentUpdate(c *gin.Context) {
	db := models.GetDB()
	comment := models.Comment{}
	if err := c.ShouldBind(&comment); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		logrus.Error(err)
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/comments/%s/edit", c.Param("id")))
		return
	}
	if err := db.Save(&comment).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/comments")
}

//CommentDelete handles POST /admin/comments/:id/delete route
func CommentDelete(c *gin.Context) {
	db := models.GetDB()
	comment := models.Comment{}
	db.First(&comment, c.Param("id"))
	if comment.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	if err := db.Delete(&comment).Error; err != nil {
		logrus.Error(err)
		c.HTML(http.StatusInternalServerError, "errors/500", gin.H{"Error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/admin/comments")
}
