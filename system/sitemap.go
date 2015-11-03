package system

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/denisbakhtin/sitemap"
	"path"
	"time"
)

func CreateXmlSitemap() {
	logrus.Info("Starting XML sitemap generation")
	folder := path.Join(GetConfig().Public, "sitemap")
	domain := "http://localhost:8080" //change in release mode
	now := time.Now()
	items := make([]sitemap.Item, 1)

	//Home page
	items = append(items, sitemap.Item{
		Loc:        fmt.Sprintf("%s", domain),
		LastMod:    now,
		Changefreq: "daily",
		Priority:   1,
	})

	//Posts
	posts, err := models.GetPublishedPosts()
	if err != nil {
		logrus.Error(err)
		return
	}
	for i := range posts {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s/posts/%d", domain, posts[i].Id),
			LastMod:    posts[i].UpdatedAt,
			Changefreq: "weekly",
			Priority:   0.9,
		})
	}

	//Static pages
	pages, err := models.GetPublishedPages()
	if err != nil {
		logrus.Error(err)
		return
	}
	for i := range pages {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s/pages/%d", domain, pages[i].Id),
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
