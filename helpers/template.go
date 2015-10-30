package helpers

import (
	"fmt"
	"time"

	"github.com/denisbakhtin/ginbasic/models"
	"github.com/gin-gonic/gin"
)

//IsActive checks uri against currently active (uri, or nil) and returns "active" if they are equal
func IsActive(active interface{}, uri string) string {
	if s, ok := active.(string); ok {
		if s == uri {
			return "active"
		}
	}
	return ""
}

//DateTime prints timestamp in human format
func DateTime(t time.Time) string {
	return fmt.Sprintf("%d-%d-%d %d:%d:%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

//DefaultH returns common to all pages template data
func DefaultH(c *gin.Context) gin.H {
	user, _ := c.Get("User")
	if u, ok := user.(*models.User); ok {
		u.Password = "" //clear password hash
	}
	signupEnabled, _ := c.Get("SignupEnabled")
	return gin.H{
		"ActiveUser":    user,          //signed in models.User
		"Active":        "",            //active uri shortening for menu item highlight
		"Title":         "",            //page title:w
		"SignupEnabled": signupEnabled, //signup route is enabled (otherwise everyone can signup ;)
	}
}
