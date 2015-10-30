package models

type Page struct {
	Id          int64  `form:"id" json:"id"`
	Name        string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description"`
	Published   bool   `form:"published" json:"published"`
}

func (page *Page) Insert() error {
	err := db.QueryRow("INSERT INTO pages(name, description, published) VALUES($1,$2,$3) RETURNING id", page.Name, page.Description, page.Published).Scan(&page.Id)
	return err
}

func (page *Page) Update() error {
	_, err := db.Exec("UPDATE pages SET name=$2, description=$3, published=$4 WHERE id=$1", page.Id, page.Name, page.Description, page.Published)
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
