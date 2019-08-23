package models

import (
	"fmt"
	"html/template"

	"github.com/microcosm-cc/bluemonday"
)

//Post type contains blog post info
type Post struct {
	Model

	Title     string `form:"title" binding:"required"`
	Content   string `form:"content"`
	Published bool   `form:"published"`
	UserID    uint64
	User      User      `binding:"-" gorm:"association_autoupdate:false;association_autocreate:false"`
	FormTags  []string  `form:"tags" gorm:"-"`
	Tags      []Tag     `binding:"-" form:"-" json:"tags" gorm:"many2many:posts_tags;"`
	Comments  []Comment `binding:"-"`
}

//Excerpt returns post excerpt by removing html tags first and truncating to 300 symbols
func (post *Post) Excerpt() template.HTML {
	//you can sanitize, cut it down, add images, etc
	policy := bluemonday.StrictPolicy() //remove all html tags
	sanitized := policy.Sanitize(post.Content)
	excerpt := template.HTML(truncate(sanitized, 300) + "...")
	return excerpt
}

//HTMLContent returns html content that won't be escaped
func (post *Post) HTMLContent() template.HTML {
	return template.HTML(post.Content)
}

//URL returns the post's canonical url
func (post *Post) URL() string {
	return fmt.Sprintf("/posts/%d", post.ID)
}
