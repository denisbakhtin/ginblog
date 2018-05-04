package models

import (
	"github.com/jinzhu/gorm"
)

//Page type contains page info
type Page struct {
	gorm.Model
	Name        string
	Description string
	Published   bool
}
