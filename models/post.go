package models

import (
	"html/template"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"gopkg.in/guregu/null.v3"
)

type Post struct {
	Id          int64     `form:"id" json:"id"`
	Name        string    `form:"name" json:"name" binding:"required"`
	Description string    `form:"description" json:"description"`
	Published   bool      `form:"published" json:"published"`
	UserId      null.Int  `form:"-" json:"user_id" db:"user_id"`
	Timestamp   time.Time `form:"-" json:"timestamp"`
	//calculated fields
	Author       User     `form:"-" json:"author" db:"author"`
	Tags         []string `form:"tags" json:"tags" db:"-"` //can't make gin Bind form field to []Tag, so use []string instead
	CommentCount int64    `form:"-" json:"comment_count" db:"comment_count"`
}

func (post *Post) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	err = db.QueryRow("INSERT INTO posts(name, description, published, user_id, timestamp) VALUES($1,$2,$3,$4,$5) RETURNING id", post.Name, post.Description, post.Published, post.UserId, time.Now()).Scan(&post.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := post.UpdateTags(tx); err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (post *Post) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE posts SET name=$2, description=$3, published=$4 WHERE id=$1", post.Id, post.Name, post.Description, post.Published)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := post.UpdateTags(tx); err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

//UpdateTags inserts new (non existent) tags & associations and removes obsolete associations
func (post *Post) UpdateTags(tx *sqlx.Tx) error {
	neu := make(map[string]bool)
	old := make(map[string]bool)
	for i := range post.Tags {
		neu[post.Tags[i]] = true
	}
	exist := make([]string, 0)
	err := db.Select(&exist, "SELECT tag_name FROM poststags WHERE post_id=$1", post.Id)
	if err != nil {
		return err
	}
	for i := range exist {
		if _, ex := neu[exist[i]]; ex {
			delete(neu, exist[i])
		} else {
			old[exist[i]] = true
		}
	}
	for name := range neu {
		//create new tag if not exists
		_, err = tx.Exec("INSERT INTO tags (name) SELECT $1 WHERE NOT EXISTS(SELECT null FROM tags WHERE name=$1)", name)
		if err != nil {
			return err
		}
		//insert new association
		_, err := tx.Exec("INSERT INTO poststags(post_id, tag_name) VALUES($1,$2)", post.Id, name)
		if err != nil {
			return err
		}
	}

	for name := range old {
		//remove association
		_, err := tx.Exec("DELETE FROM poststags WHERE post_id=$1 AND tag_name=$2", post.Id, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (post *Post) Delete() error {
	_, err := db.Exec("DELETE FROM posts WHERE id=$1", post.Id)
	return err
}

func (post *Post) Excerpt() template.HTML {
	//you can sanitize, cut it down, add images, etc
	policy := bluemonday.StrictPolicy() //remove all html tags
	sanitized := policy.Sanitize(string(blackfriday.MarkdownCommon([]byte(post.Description))))
	excerpt := template.HTML(truncate(sanitized, 300) + "...")
	return excerpt
}

func GetPost(id interface{}) (*Post, error) {
	post := &Post{}
	err := db.Get(post, "SELECT * FROM posts WHERE id=$1", id)
	if err != nil {
		return post, err
	}
	err = db.Get(&post.Author, "SELECT id,name FROM users WHERE id=$1", post.UserId)
	if err != nil {
		return post, err
	}
	err = db.Select(&post.Tags, "SELECT name FROM tags WHERE EXISTS (SELECT null FROM poststags WHERE post_id=$1 AND tag_name=tags.name)", id)
	return post, err
}

func GetPosts() ([]Post, error) {
	var list []Post
	err := db.Select(&list, "SELECT * FROM posts ORDER BY posts.id DESC")
	return list, err
}

func GetRecentPosts() ([]Post, error) {
	var list []Post
	err := db.Select(&list, "SELECT id, name FROM posts WHERE published=$1 ORDER BY id DESC LIMIT 7", true)
	return list, err
}

func GetPostMonths() ([]Post, error) {
	var list []Post
	err := db.Select(&list, "SELECT DISTINCT date_trunc('month', timestamp) as timestamp FROM posts WHERE published=$1 ORDER BY timestamp DESC", true)
	return list, err
}

func GetPostsByArchive(year, month int) ([]Post, error) {
	var list []Post
	err := db.Select(&list, "SELECT * FROM posts WHERE published=$1 AND date_part('year', timestamp)=$2 AND date_part('month', timestamp)=$3 ORDER BY timestamp DESC", true, year, month)
	if err != nil {
		return list, err
	}
	for i := range list {
		err := db.Get(&list[i].Author, "SELECT id,name FROM users WHERE id=$1", list[i].UserId)
		if err != nil {
			return list, err
		}
		err = db.Select(&list[i].Tags, "SELECT name FROM tags WHERE EXISTS (SELECT null FROM poststags WHERE post_id=$1 AND tag_name=tags.name)", list[i].Id)
		if err != nil {
			return list, err
		}
	}
	return list, err
}

func GetPostsByTag(name string) ([]Post, error) {
	var list []Post
	err := db.Select(&list, "SELECT * FROM posts WHERE published=$1 AND EXISTS (SELECT null FROM poststags WHERE poststags.post_id=posts.id AND poststags.tag_name=$2) ORDER BY timestamp DESC", true, name)
	if err != nil {
		return list, err
	}
	for i := range list {
		err := db.Get(&list[i].Author, "SELECT id,name FROM users WHERE id=$1", list[i].UserId)
		if err != nil {
			return list, err
		}
		err = db.Select(&list[i].Tags, "SELECT name FROM tags WHERE EXISTS (SELECT null FROM poststags WHERE post_id=$1 AND tag_name=tags.name)", list[i].Id)
		if err != nil {
			return list, err
		}
	}
	return list, err
}
