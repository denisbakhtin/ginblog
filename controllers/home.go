package controllers

import (
	"net/http"

	"github.com/denisbakhtin/ginbasic/helpers"
	"github.com/gin-gonic/gin"
)

func HomeGet(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "Welcome to basic GIN web-site"
	h["Active"] = "home"
	c.HTML(http.StatusOK, "home/show", h)
}
