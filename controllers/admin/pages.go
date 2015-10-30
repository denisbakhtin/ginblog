package admin

import (
	"net/http"

	"github.com/denisbakhtin/ginbasic/helpers"
	"github.com/denisbakhtin/ginbasic/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GET page list
func PageIndex(c *gin.Context) {
	list, err := models.GetPageList()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "List of pages"
	h["List"] = list
	h["Active"] = "pages"
	c.HTML(http.StatusOK, "admin/pages/index", h)
}

// GET page creation form
func PageNew(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "New page"
	h["Active"] = "pages"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "admin/pages/form", h)
}

// POST page creation form
func PageCreate(c *gin.Context) {
	page := &models.Page{}
	if err := c.Bind(page); err == nil {
		if err := page.Insert(); err != nil {
			c.HTML(http.StatusInternalServerError, "errors/500", nil)
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
	c.HTML(http.StatusOK, "admin/pages/form", h)
}

// POST page update form
func PageUpdate(c *gin.Context) {
	page := &models.Page{}
	if err := c.Bind(page); err == nil {
		if err := page.Update(); err != nil {
			c.HTML(http.StatusInternalServerError, "errors/500", nil)
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
		return
	} else {
		c.Redirect(http.StatusFound, "/admin/pages")
	}
}
