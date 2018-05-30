package helpers

import (
	"github.com/gin-gonic/gin"
)

//DefaultH returns common to all pages template data
func DefaultH(c *gin.Context) gin.H {
	return gin.H{
		"Title":   "", //page title:w
		"Context": c,
	}
}
