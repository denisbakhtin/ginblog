package controllers

import (
	"net/http"

	"html/template"

	"github.com/denisbakhtin/ginbasic/helpers"
	"github.com/denisbakhtin/ginbasic/models"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
)

// GET /pages/:id route
func PageGet(c *gin.Context) {
	page, err := models.GetPage(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		return
	}
	if !page.Published {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = page.Name
	h["Description"] = template.HTML(string(blackfriday.MarkdownCommon([]byte(page.Description))))
	h["Active"] = "pages"
	c.HTML(http.StatusOK, "pages/show", h)
}
