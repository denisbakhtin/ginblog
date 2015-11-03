package models

import (
	"time"
)

//Page type contains page info
type Page struct {
	ID          int64     `form:"id" json:"id" database:"id"`
	Name        string    `form:"name" json:"name"`
	Description string    `form:"description" json:"description"`
	Published   bool      `form:"published" json:"published"`
	CreatedAt   time.Time `form:"_" json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `form:"_" json:"updated_at" db:"updated_at"`
}

//Insert saves Page struct into db
func (page *Page) Insert() error {
	err := db.QueryRow("INSERT INTO pages(name, description, published, created_at, updated_at) VALUES($1,$2,$3,$4,$4) RETURNING id", page.Name, page.Description, page.Published, time.Now()).Scan(&page.ID)
	return err
}

//Update saves Page changes into db
func (page *Page) Update() error {
	_, err := db.Exec("UPDATE pages SET name=$2, description=$3, published=$4, updated_at=$5 WHERE id=$1", page.ID, page.Name, page.Description, page.Published, time.Now())
	return err
}

//Delete removes page from db
func (page *Page) Delete() error {
	_, err := db.Exec("DELETE FROM pages WHERE id=$1", page.ID)
	return err
}

//GetPage fetches page from db by its id
func GetPage(id interface{}) (*Page, error) {
	page := &Page{}
	err := db.Get(page, "SELECT * FROM pages WHERE id=$1", id)
	return page, err
}

//GetPages returns a slice of all pages
func GetPages() ([]Page, error) {
	var list []Page
	err := db.Select(&list, "SELECT * FROM pages ORDER BY id")
	return list, err
}

//GetPublishedPages returns a slice of pages with .Published=true
func GetPublishedPages() ([]Page, error) {
	var list []Page
	err := db.Select(&list, "SELECT * FROM pages WHERE published=$1 ORDER BY id", true)
	return list, err
}
