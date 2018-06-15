package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/gin"
)

//ArchiveGet handles GET /archives/:year/:month route
func ArchiveGet(c *gin.Context) {
	year, _ := strconv.Atoi(c.Param("year"))
	month, _ := strconv.Atoi(c.Param("month"))
	db := models.GetDB()
	var posts []models.Post
	db.Where("published = true AND date_part('year', created_at)=? AND date_part('month', created_at)=?", year, month).Order("created_at desc").Find(&posts)
	h := DefaultH(c)
	h["Title"] = fmt.Sprintf("%s %d archives", time.Month(month).String(), year)
	h["Posts"] = posts
	c.HTML(http.StatusOK, "archives/show", h)
}
