package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//PageGet handles GET /pages/:id route
func PageGet(c *gin.Context) {
	db := models.GetDB()
	page := models.Page{}
	db.First(&page, c.Param("id"))
	if page.ID == 0 || !page.Published {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = page.Title
	h["Page"] = page
	c.HTML(http.StatusOK, "pages/show", h)
}

//PageIndex handles GET /admin/pages route
func PageIndex(c *gin.Context) {
	db := models.GetDB()
	var pages []models.Page
	db.Find(&pages)
	h := DefaultH(c)
	h["Title"] = "List of pages"
	h["Pages"] = pages
	c.HTML(http.StatusOK, "pages/index", h)
}

//PageNew handles GET /admin/new_page route
func PageNew(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "New page"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "pages/form", h)
}

//PageCreate handles POST /admin/new_page route
func PageCreate(c *gin.Context) {
	db := models.GetDB()
	page := models.Page{}
	if err := c.ShouldBind(&page); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_page")
		return
	}

	if err := db.Create(&page).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/pages")
}

//PageEdit handles GET /admin/pages/:id/edit route
func PageEdit(c *gin.Context) {
	db := models.GetDB()
	page := models.Page{}
	db.First(&page, c.Param("id"))
	if page.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = "Edit page"
	h["Page"] = page
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "pages/form", h)
}

//PageUpdate handles POST /admin/pages/:id/edit route
func PageUpdate(c *gin.Context) {
	page := &models.Page{}
	db := models.GetDB()
	if err := c.ShouldBind(page); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/pages")
		return
	}
	if err := db.Save(&page).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/pages")
}

//PageDelete handles POST /admin/pages/:id/delete route
func PageDelete(c *gin.Context) {
	page := models.Page{}
	db := models.GetDB()
	db.First(&page, c.Param("id"))
	if page.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	if err := db.Delete(&page).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/pages")
}
