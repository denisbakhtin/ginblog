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
	list, err := models.GetPostsByArchive(year, month)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = fmt.Sprintf("%s %d archives", time.Month(month).String(), year)
	h["List"] = list
	h["Active"] = fmt.Sprintf("archives/%d/%02d", year, month)
	c.HTML(http.StatusOK, "archives/show", h)
}
