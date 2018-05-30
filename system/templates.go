package system

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"

	"github.com/denisbakhtin/ginblog/models"
)

var tmpl *template.Template

//LoadTemplates loads templates from views directory
func LoadTemplates() {
	tmpl = template.New("").Funcs(template.FuncMap{
		"isActiveLink":        isActiveLink,
		"stringInSlice":       stringInSlice,
		"formatDateTime":      formatDateTime,
		"recentPosts":         recentPosts,
		"tags":                tags,
		"archives":            archives,
		"now":                 now,
		"activeUserEmail":     activeUserEmail,
		"activeUserID":        activeUserID,
		"isUserAuthenticated": isUserAuthenticated,
		"signUpEnabled":       signUpEnabled,
		"csrfToken":           csrfToken,
	})
	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".gohtml") {
			var err error
			tmpl, err = tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if err := filepath.Walk("views", fn); err != nil {
		panic(err)
	}
}

//GetTemplates returns preloaded templates
func GetTemplates() *template.Template {
	return tmpl
}

//isActiveLink checks uri against currently active (uri, or nil) and returns "active" if they are equal
func isActiveLink(c *gin.Context, uri string) string {
	if c.Request.RequestURI == uri {
		return "active"
	}
	return ""
}

//formatDateTime prints timestamp in human format
func formatDateTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

//stringInSlice returns true if value is in list slice
func stringInSlice(value string, list []string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

//recentPosts returns the list of recent blog posts
func recentPosts() []models.Post {
	db := models.GetDB()
	var list []models.Post
	db.Where("published = true").Order("id desc").Limit(7).Find(&list)
	return list
}

//tags returns the list of blog tags
func tags() []models.Tag {
	var list []models.Tag
	//list, _ := models.GetNotEmptyTags()
	return list
}

//archives returns the list of blog archives
func archives() []models.Post {
	db := models.GetDB()
	var list []models.Post
	db.Select("distinct date_trunc('month', created_at) as created_at").Where("published = true").Order("created_at desc").Find(&list)
	return list
}

//now returns current timestamp
func now() time.Time {
	return time.Now()
}

//activeUserEmail returns current authenticated user email
func activeUserEmail(c *gin.Context) string {
	u, _ := c.Get("User")
	if user, ok := u.(*models.User); ok {
		return user.Email
	}
	return ""
}

//activeUserID returns current authenticated user ID
func activeUserID(c *gin.Context) uint {
	u, _ := c.Get("User")
	if user, ok := u.(*models.User); ok {
		return user.ID
	}
	return 0
}

//isUserAuthenticated returns true is user is authenticated
func isUserAuthenticated(c *gin.Context) bool {
	u, _ := c.Get("User")
	if _, ok := u.(*models.User); ok {
		return true
	}
	return false
}

//signUpEnabled returns true if sign up is enabled by config
func signUpEnabled(c *gin.Context) bool {
	se, _ := c.Get("SignupEnabled")
	if enabled, ok := se.(bool); ok {
		return enabled
	}
	return false
}

//csrfToken returns CSRF protection token
func csrfToken(c *gin.Context) string {
	return csrf.GetToken(c)
}
