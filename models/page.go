package models

import (
	"time"
)

type Page struct {
	Id          int64     `form:"id" json:"id"`
	Name        string    `form:"name" json:"name"`
	Description string    `form:"description" json:"description"`
	Published   bool      `form:"published" json:"published"`
	CreatedAt   time.Time `form:"_" json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `form:"_" json:"updated_at" db:"updated_at"`
}

func (page *Page) Insert() error {
	err := db.QueryRow("INSERT INTO pages(name, description, published, created_at, updated_at) VALUES($1,$2,$3,$4,$4) RETURNING id", page.Name, page.Description, page.Published, time.Now()).Scan(&page.Id)
	return err
}

func (page *Page) Update() error {
	_, err := db.Exec("UPDATE pages SET name=$2, description=$3, published=$4, updated_at=$5 WHERE id=$1", page.Id, page.Name, page.Description, page.Published, time.Now())
	return err
}

func (page *Page) Delete() error {
	_, err := db.Exec("DELETE FROM pages WHERE id=$1", page.Id)
	return err
}

func GetPage(id interface{}) (*Page, error) {
	page := &Page{}
	err := db.Get(page, "SELECT * FROM pages WHERE id=$1", id)
	return page, err
}

func GetPages() ([]Page, error) {
	var list []Page
	err := db.Select(&list, "SELECT * FROM pages ORDER BY id")
	return list, err
}

func GetPublishedPages() ([]Page, error) {
	var list []Page
	err := db.Select(&list, "SELECT * FROM pages WHERE published=$1 ORDER BY id", true)
	return list, err
}
