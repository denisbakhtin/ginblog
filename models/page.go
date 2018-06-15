package models

//Page type contains page info
type Page struct {
	Model

	Title     string `form:"title"`
	Content   string `form:"content"`
	Published bool   `form:"published"`
}
