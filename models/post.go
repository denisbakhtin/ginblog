package models

import (
	"html/template"

	"github.com/jinzhu/gorm"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

//Post type contains blog post info
type Post struct {
	gorm.Model
	Name        string
	Description string
	Published   bool
	UserID      uint
	User        User
	//calculated fields
	Tags         []string `form:"tags" json:"tags" db:"-"` //can't make gin Bind form field to []Tag, so use []string instead
	CommentCount int64    `form:"-" json:"comment_count" db:"comment_count"`
}

//Excerpt returns post excerpt by removing html tags first and truncating to 300 symbols
func (post *Post) Excerpt() template.HTML {
	//you can sanitize, cut it down, add images, etc
	policy := bluemonday.StrictPolicy() //remove all html tags
	sanitized := policy.Sanitize(string(blackfriday.MarkdownCommon([]byte(post.Description))))
	excerpt := template.HTML(truncate(sanitized, 300) + "...")
	return excerpt
}
