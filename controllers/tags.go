package controllers

import (
	"net/http"

	"github.com/denisbakhtin/ginblog/helpers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/gin"
)

// GET /tags/:name route
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
