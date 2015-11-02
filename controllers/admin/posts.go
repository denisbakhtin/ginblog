package admin

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/helpers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v3"
)

// GET post entry list
func PostIndex(c *gin.Context) {
	list, err := models.GetPosts()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "List of blog posts"
	h["List"] = list
	h["Active"] = "posts"
	c.HTML(http.StatusOK, "admin/posts/index", h)
}

// GET post creation form
func PostNew(c *gin.Context) {
	tags, _ := models.GetTags()
	h := helpers.DefaultH(c)
	h["Title"] = "New post entry"
	h["Active"] = "posts"
	h["Tags"] = tags
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "admin/posts/form", h)
}

// POST post creation form
func PostCreate(c *gin.Context) {
	post := &models.Post{}
	if err := c.Bind(post); err == nil {
		if user, exists := c.Get("User"); exists {
			post.UserId = null.IntFrom(user.(*models.User).Id)
		}
		if err := post.Insert(); err != nil {
			c.HTML(http.StatusInternalServerError, "errors/500", nil)
			return
		}
		c.Redirect(http.StatusFound, "/admin/posts")
	} else {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_post")
	}
}

// GET post update form
func PostEdit(c *gin.Context) {
	post, _ := models.GetPost(c.Param("id"))
	if post.Id == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	tags, _ := models.GetTags()
	h := helpers.DefaultH(c)
	h["Title"] = "Edit post entry"
	h["Active"] = "posts"
	h["Post"] = post
	h["Tags"] = tags
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "admin/posts/form", h)
}

// POST post update form
func PostUpdate(c *gin.Context) {
	post := &models.Post{}
	if err := c.Bind(post); err == nil {
		logrus.Warn(post)
		if err := post.Update(); err != nil {
			logrus.Warn(err)
			c.HTML(http.StatusInternalServerError, "errors/500", nil)
			return
		}
		c.Redirect(http.StatusFound, "/admin/posts")
	} else {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/posts/%s/edit", c.Param("id")))
	}
}

// POST post deletion request
func PostDelete(c *gin.Context) {
	post, _ := models.GetPost(c.Param("id"))
	if err := post.Delete(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		return
	} else {
		c.Redirect(http.StatusFound, "/admin/posts")
	}
}
