package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//AdminGet handles GET /admin/ route
func AdminGet(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "Admin dashboard"
	c.HTML(http.StatusOK, "dashboard/show", h)
}
