package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//NotFound handles gin NotFound error
func NotFound(c *gin.Context) {
	ShowErrorPage(c, http.StatusNotFound, nil)
}

//MethodNotAllowed handles gin MethodNotAllowed error
func MethodNotAllowed(c *gin.Context) {
	ShowErrorPage(c, http.StatusMethodNotAllowed, nil)
}

func ShowErrorPage(c *gin.Context, code int, err error) {
	H := DefaultH(c)
	H["Error"] = err
	c.HTML(code, fmt.Sprintf("errors/%d", code), H)
}
