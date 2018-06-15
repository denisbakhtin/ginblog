package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//HomeGet handles GET / route
func HomeGet(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "Welcome to basic GIN blog"
	c.HTML(http.StatusOK, "home/show", h)
}
