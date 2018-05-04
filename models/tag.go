package models

import (
	"github.com/jinzhu/gorm"
)

//Tag struct contains post tag info
type Tag struct {
	gorm.Model
	Name string
	//calculated fields
	PostCount int64 `form:"-"`
}
