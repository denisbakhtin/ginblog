package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/denisbakhtin/ginblog/models"
	"github.com/denisbakhtin/ginblog/system"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//ContextData stores in gin context the common data, such as user info...
func ContextData() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if uID := session.Get(userIDKey); uID != nil {
			user := models.User{}
			models.GetDB().First(&user, uID)
			if user.ID != 0 {
				c.Set("User", &user)
			}
		}

		if system.GetConfig().SignupEnabled {
			c.Set("SignupEnabled", true)
		}
		c.Next()
	}
}

//AuthRequired grants access to authenticated users, requires SharedData middleware
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("User"); user != nil {
			c.Next()
		} else {
			c.Redirect(http.StatusFound, fmt.Sprintf("/signin?return=%s", url.QueryEscape(c.Request.RequestURI)))
			c.Abort()
		}
	}
}
