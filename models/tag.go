package models

import (
	"github.com/jinzhu/gorm"
)

//Tag struct contains post tag info
type Tag struct {
	gorm.Model
	Name string
	Posts []Post `gorm:"many2many:posts_tags;"`
	PostCount int64 `form:"-"`
}
