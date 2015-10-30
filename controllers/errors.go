package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "errors/404", nil)
}
func MethodNotAllowed(c *gin.Context) {
	c.HTML(http.StatusMethodNotAllowed, "errors/405", nil)
}
