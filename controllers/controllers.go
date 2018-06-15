package controllers

import (
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

const userIDKey = "UserID"

//DefaultH returns common to all pages template data
func DefaultH(c *gin.Context) gin.H {
	return gin.H{
		"Title":   "", //page title:w
		"Context": c,
		"Csrf":    csrf.GetToken(c),
	}
}
