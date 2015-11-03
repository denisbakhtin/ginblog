package controllers

import (
	"net/http"

	"html/template"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/helpers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
)

// GET /pages/:id route
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

// GET page list
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

// GET page creation form
func PageNew(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "New page"
	h["Active"] = "pages"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "pages/form", h)
}

// POST page creation form
func PageCreate(c *gin.Context) {
	page := &models.Page{}
	if err := c.Bind(page); err == nil {
		if err := page.Insert(); err != nil {
			c.HTML(http.StatusInternalServerError, "errors/500", nil)
			logrus.Error(err)
			return
		}
		c.Redirect(http.StatusFound, "/admin/pages")
	} else {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_page")
	}
}

// GET page update form
func PageEdit(c *gin.Context) {
	page, _ := models.GetPage(c.Param("id"))
	if page.Id == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "Edit page"
	h["Active"] = "pages"
	h["Page"] = page
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "pages/form", h)
}

// POST page update form
func PageUpdate(c *gin.Context) {
	page := &models.Page{}
	if err := c.Bind(page); err == nil {
		if err := page.Update(); err != nil {
			c.HTML(http.StatusInternalServerError, "errors/500", nil)
			logrus.Error(err)
			return
		}
		c.Redirect(http.StatusFound, "/admin/pages")
	} else {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/pages")
	}
}

// POST page deletion request
func PageDelete(c *gin.Context) {
	page, _ := models.GetPage(c.Param("id"))
	if err := page.Delete(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	} else {
		c.Redirect(http.StatusFound, "/admin/pages")
	}
}
