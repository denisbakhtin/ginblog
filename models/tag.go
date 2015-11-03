package models

import ()

//Tag struct contains post tag info
type Tag struct {
	Name string `form:"name" json:"name"`
	//calculated fields
	PostCount int64 `form:"-" json:"post_count" db:"post_count"`
}

//Insert saves tag info into db
func (tag *Tag) Insert() error {
	_, err := db.Exec("INSERT INTO tags(name) VALUES($1)", tag.Name)
	return err
}

//Delete removes tag from db, according to postgresql constraints tag associations are removed either
func (tag *Tag) Delete() error {
	_, err := db.Exec("DELETE FROM tags WHERE name=$1", tag.Name)
	return err
}

//GetTag returns tag by its name
func GetTag(name interface{}) (*Tag, error) {
	tag := &Tag{}
	err := db.Get(tag, "SELECT * FROM tags WHERE name=$1", name)
	return tag, err
}

//GetTags returns a slice of tags, ordered by name
func GetTags() ([]Tag, error) {
	var list []Tag
	err := db.Select(&list, "SELECT * FROM tags ORDER BY name ASC")
	return list, err
}

//GetNotEmptyTags returns a slice of tags that have at least one associated blog post
func GetNotEmptyTags() ([]Tag, error) {
	var list []Tag
	err := db.Select(&list, "SELECT * FROM tags WHERE EXISTS (SELECT null FROM poststags WHERE poststags.tag_name=tags.name) ORDER BY name ASC")
	return list, err
}
