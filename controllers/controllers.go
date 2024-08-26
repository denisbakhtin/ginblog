package controllers

import (
	"fmt"
	"log/slog"
	"path"
	"time"

	"github.com/denisbakhtin/ginblog/config"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/denisbakhtin/sitemap"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

const userIDKey = "UserID"

// DefaultH returns common to all pages template data
func DefaultH(c *gin.Context) gin.H {
	return gin.H{
		"Title":   "", //page title:w
		"Context": c,
		"Csrf":    csrf.GetToken(c),
	}
}

// CreateXMLSitemap creates xml sitemap for search engines, and saves in public/sitemap folder
func CreateXMLSitemap() {
	slog.Info("Starting XML sitemap generation")
	folder := path.Join(config.GetConfig().Public, "sitemap")
	domain := config.GetConfig().Domain
	now := time.Now()
	items := make([]sitemap.Page, 1)
	db := models.GetDB()

	//Home page
	items = append(items, sitemap.Page{
		Loc:        domain,
		LastMod:    now,
		Changefreq: sitemap.Daily,
		Priority:   1,
	})

	//Posts
	var posts []models.Post
	db.Where("published = true").Find(&posts)
	for i := range posts {
		items = append(items, sitemap.Page{
			Loc:        fmt.Sprintf("%s%s", domain, posts[i].URL()),
			LastMod:    posts[i].UpdatedAt,
			Changefreq: sitemap.Weekly,
			Priority:   0.9,
		})
	}

	//Static pages
	var pages []models.Page
	db.Where("published = true").Find(&pages)
	for i := range pages {
		items = append(items, sitemap.Page{
			Loc:        fmt.Sprintf("%s%s", domain, pages[i].URL()),
			LastMod:    pages[i].UpdatedAt,
			Changefreq: sitemap.Monthly,
			Priority:   0.8,
		})
	}
	if err := sitemap.SiteMap(path.Join(folder, "sitemap1.xml.gz"), items); err != nil {
		slog.Error(err.Error())
		return
	}
	if err := sitemap.SiteMapIndex(folder, "sitemap_index.xml", domain+"/public/sitemap/"); err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("XML sitemap has been generated in " + folder)
}
