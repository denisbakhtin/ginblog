package models

import "fmt"

//Page type contains page info
type Page struct {
	Model

	Title     string `form:"title"`
	Content   string `form:"content"`
	Published bool   `form:"published"`
}

//URL returns the page's canonical url
func (page *Page) URL() string {
	return fmt.Sprintf("/pages/%d", page.ID)
}
