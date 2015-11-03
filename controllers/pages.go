package controllers

import (
	"net/http"

	"html/template"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/helpers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
)

//PageGet handles GET /pages/:id route
func PageGet(c *gin.Context) {
	page, err := models.GetPage(c.Param("id"))
	if err != nil || !page.Published {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = page.Name
	h["Description"] = template.HTML(string(blackfriday.MarkdownCommon([]byte(page.Description))))
	h["Active"] = "pages"
	c.HTML(http.StatusOK, "pages/show", h)
}

//PageIndex handles GET /admin/pages route
func PageIndex(c *gin.Context) {
	list, err := models.GetPages()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "List of pages"
	h["List"] = list
	h["Active"] = "pages"
	c.HTML(http.StatusOK, "pages/index", h)
}

//PageNew handles GET /admin/new_page route
func PageNew(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "New page"
	h["Active"] = "pages"
	session := helpers.GetSession(c)
	h["Flash"] = session.Flashes()
	session.Save(c.Request, c.Writer)

	c.HTML(http.StatusOK, "pages/form", h)
}

//PageCreate handles POST /admin/new_page route
func PageCreate(c *gin.Context) {
	page := &models.Page{}
	if err := c.Bind(page); err != nil {
		session := helpers.GetSession(c)
		session.AddFlash(err.Error())
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusSeeOther, "/admin/new_page")
		return
	}

	if err := page.Insert(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/pages")
}

//PageEdit handles GET /admin/pages/:id/edit route
func PageEdit(c *gin.Context) {
	page, _ := models.GetPage(c.Param("id"))
	if page.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "Edit page"
	h["Active"] = "pages"
	h["Page"] = page
	session := helpers.GetSession(c)
	h["Flash"] = session.Flashes()
	session.Save(c.Request, c.Writer)
	c.HTML(http.StatusOK, "pages/form", h)
}

//PageUpdate handles POST /admin/pages/:id/edit route
func PageUpdate(c *gin.Context) {
	page := &models.Page{}
	if err := c.Bind(page); err != nil {
		session := helpers.GetSession(c)
		session.AddFlash(err.Error())
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusSeeOther, "/admin/pages")
		return
	}
	if err := page.Update(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/pages")
}

//PageDelete handles POST /admin/pages/:id/delete route
func PageDelete(c *gin.Context) {
	page, _ := models.GetPage(c.Param("id"))
	if err := page.Delete(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/pages")
}
