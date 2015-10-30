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

// GET /archives/:year/:month route
func ArchiveGet(c *gin.Context) {
	year, _ := strconv.Atoi(c.Param("year"))
	month, _ := strconv.Atoi(c.Param("month"))
	list, err := models.GetPostsByYearMonth(year, month)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = fmt.Sprintf("%s %d archives", helpers.MonthHuman(time.Month(month)), year)
	h["List"] = list
	h["Active"] = "archives"
	c.HTML(http.StatusOK, "archives/show", h)
}
