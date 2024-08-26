package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Tag struct contains post tag info
type Tag struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	Title string `binding:"required" form:"title" gorm:"primary_key"`
	Slug  string `binding:"-"`
	Posts []Post `gorm:"many2many:posts_tags;foreignkey:title"`
}

// BeforeSave gorm hook
func (tag *Tag) BeforeSave(tx *gorm.DB) (err error) {
	tag.Slug = createSlug(tag.Title)
	return
}

// BeforeDelete gorm hook removes a join record
func (tag *Tag) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Exec("DELETE FROM posts_tags WHERE tag_title = ?", tag.Title).Error
}

// URL returns the tag's canonical url
func (tag *Tag) URL() string {
	return fmt.Sprintf("/tags/%s", tag.Slug)
}
