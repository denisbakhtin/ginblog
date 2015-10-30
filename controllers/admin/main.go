package admin

import (
	"net/http"

	"github.com/denisbakhtin/ginbasic/helpers"
	"github.com/gin-gonic/gin"
)

func AdminGet(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "Admin dashboard"
	c.HTML(http.StatusOK, "admin/home/show", h)
}
