package controllers

import (
	"fmt"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"time"
)

// GET /rss route
func RssGet(c *gin.Context) {
	now := time.Now()

	feed := &feeds.Feed{
		Title:       "ginblog",
		Link:        &feeds.Link{Href: "http://localhost:8080"},
		Description: "GIN-powered blog boilerplate",
		Author:      &feeds.Author{"Blog Author", "author@blog.net"}, //hide email from spammers?
		Created:     now,
		Copyright:   "This work is copyright © Ginblog",
	}

	feed.Items = make([]*feeds.Item, 0)
	posts, err := models.GetPublishedPosts()
	if err != nil {
		c.HTML(404, "errors/404", nil)
		return
	}

	for i := range posts {
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          fmt.Sprintf("http://localhost:8080/posts/%d", posts[i].Id),
			Title:       posts[i].Name,
			Link:        &feeds.Link{Href: fmt.Sprintf("http://localhost:8080/posts/%d", posts[i].Id)},
			Description: string(posts[i].Excerpt()),
			Author:      &feeds.Author{Name: posts[i].Author.Name},
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
