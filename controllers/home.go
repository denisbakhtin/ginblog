package controllers

import (
	"net/http"

	"github.com/denisbakhtin/ginblog/helpers"
	"github.com/gin-gonic/gin"
)

func HomeGet(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "Welcome to basic GIN blog"
	h["Active"] = "home"
	c.HTML(http.StatusOK, "home/show", h)
}
