package models

import ()

type Tag struct {
	Name string `form:"name" json:"name" binding:"required"`
	//calculated fields
	PostCount int64 `form:"-" json:"post_count" db:"post_count"`
}

func (tag *Tag) Insert() error {
	_, err := db.Exec("INSERT INTO tags(name) VALUES($1)", tag.Name)
	return err
}

func (tag *Tag) Delete() error {
	_, err := db.Exec("DELETE FROM tags WHERE name=$1", tag.Name)
	return err
}

func GetTag(name interface{}) (*Tag, error) {
	tag := &Tag{}
	err := db.Get(tag, "SELECT * FROM tags WHERE name=$1", name)
	return tag, err
}

func GetTags() ([]Tag, error) {
	var list []Tag
	err := db.Select(&list, "SELECT * FROM tags ORDER BY name ASC")
	return list, err
}

func GetNotEmptyTags() ([]Tag, error) {
	var list []Tag
	err := db.Select(&list, "SELECT * FROM tags WHERE EXISTS (SELECT null FROM poststags WHERE poststags.tag_name=tags.name) ORDER BY name ASC")
	return list, err
}
