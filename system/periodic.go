package system

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/denisbakhtin/sitemap"
	"path"
	"time"
)

//CreateXMLSitemap creates xml sitemap for search engines, and saves in public/sitemap folder
func CreateXMLSitemap() {
	logrus.Info("Starting XML sitemap generation")
	folder := path.Join(GetConfig().Public, "sitemap")
	domain := GetConfig().Domain
	now := time.Now()
	items := make([]sitemap.Item, 1)
	db := models.GetDB()

	//Home page
	items = append(items, sitemap.Item{
		Loc:        fmt.Sprintf("%s", domain),
		LastMod:    now,
		Changefreq: "daily",
		Priority:   1,
	})

	//Posts
	var posts []models.Post
	db.Where("published = true").Find(&posts)
	for i := range posts {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s/posts/%d", domain, posts[i].ID),
			LastMod:    posts[i].UpdatedAt,
			Changefreq: "weekly",
			Priority:   0.9,
		})
	}

	//Static pages
	var pages []models.Page
	db.Where("published = true").Find(&pages)
	for i := range pages {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s/pages/%d", domain, pages[i].ID),
			LastMod:    pages[i].UpdatedAt,
			Changefreq: "monthly",
			Priority:   0.8,
		})
	}
	if err := sitemap.SiteMap(path.Join(folder, "sitemap1.xml.gz"), items); err != nil {
		logrus.Error(err)
		return
	}
	if err := sitemap.SiteMapIndex(folder, "sitemap_index.xml", domain+"/public/sitemap/"); err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("XML sitemap has been generated in " + folder)
}
