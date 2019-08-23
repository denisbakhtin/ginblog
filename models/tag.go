package models

import (
	"fmt"
	"time"
)

//Tag struct contains post tag info
type Tag struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	Title string `binding:"required" form:"title" gorm:"primary_key"`
	Slug  string `binding:"-"`
	Posts []Post `gorm:"many2many:posts_tags;foreignkey:title"`
}

//BeforeSave gorm hook
func (tag *Tag) BeforeSave() (err error) {
	tag.Slug = createSlug(tag.Title)
	return
}

//URL returns the tag's canonical url
func (tag *Tag) URL() string {
	return fmt.Sprintf("/tags/%s", tag.Slug)
}
