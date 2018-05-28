package system

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/denisbakhtin/ginblog/helpers"
)

var tmpl *template.Template

//LoadTemplates loads templates from views directory
func LoadTemplates() {
	tmpl = template.New("").Funcs(template.FuncMap{
		"isActive":      helpers.IsActive,
		"stringInSlice": helpers.StringInSlice,
		"dateTime":      helpers.DateTime,
		"recentPosts":   helpers.RecentPosts,
		"tags":          helpers.Tags,
		"archives":      helpers.Archives,
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
