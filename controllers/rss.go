package controllers

import (
	"fmt"
	"time"

	"github.com/denisbakhtin/ginblog/config"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

//RssGet handles GET /rss route
func RssGet(c *gin.Context) {
	now := time.Now()
	domain := config.GetConfig().Domain
	db := models.GetDB()

	feed := &feeds.Feed{
		Title:       "ginblog",
		Link:        &feeds.Link{Href: domain},
		Description: "GIN-powered blog boilerplate",
		Author:      &feeds.Author{Name: "Blog Author", Email: "author@blog.net"}, //hide email from spammers?
		Created:     now,
		Copyright:   "This work is copyright © Ginblog",
	}

	feed.Items = make([]*feeds.Item, 0)
	var posts []models.Post
	db.Where("published = true").Find(&posts)

	for i := range posts {
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          fmt.Sprintf("%s%s", domain, posts[i].URL()),
			Title:       posts[i].Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s%s", domain, posts[i].URL())},
			Description: string(posts[i].Excerpt()),
			Author:      &feeds.Author{Name: posts[i].User.Name},
			Created:     now,
		})
	}

	rss, err := feed.ToRss()
	if err != nil {
		c.HTML(500, "errors/500", nil)
		return
	}
	c.Data(200, "text/xml", []byte(rss))
}
