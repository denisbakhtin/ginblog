package controllers

import (
	"fmt"
	"net/http"

	"html/template"

	"github.com/denisbakhtin/ginblog/helpers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
)

// GET /posts/:id route
func PostGet(c *gin.Context) {
	post, err := models.GetPost(c.Param("id"))
	if err != nil || !post.Published {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = post.Name
	h["Description"] = template.HTML(string(blackfriday.MarkdownCommon([]byte(post.Description))))
	h["Post"] = post
	h["Active"] = fmt.Sprintf("posts/%d", post.Id)
	c.HTML(http.StatusOK, "posts/show", h)
}
