package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/denisbakhtin/ginblog/helpers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/gin"
)

//ArchiveGet handles GET /archives/:year/:month route
func ArchiveGet(c *gin.Context) {
	year, _ := strconv.Atoi(c.Param("year"))
	month, _ := strconv.Atoi(c.Param("month"))
	db := models.GetDB()
	var list []models.Post
	db.Where("published = true AND date_part('year', created_at)=$1 AND date_part('month', created_at)=$2", year, month).Order("created_at desc").Find(&list)
	h := helpers.DefaultH(c)
	h["Title"] = fmt.Sprintf("%s %d archives", time.Month(month).String(), year)
	h["List"] = list
	h["Active"] = fmt.Sprintf("archives/%d/%02d", year, month)
	c.HTML(http.StatusOK, "archives/show", h)
}
