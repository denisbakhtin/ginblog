package models

import (
	"html/template"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"gopkg.in/guregu/null.v3"
)

//Post type contains blog post info
type Post struct {
	ID          int64     `form:"id" json:"id" database:"id"`
	Name        string    `form:"name" json:"name"`
	Description string    `form:"description" json:"description"`
	Published   bool      `form:"published" json:"published"`
	UserID      null.Int  `form:"-" json:"user_id" db:"user_id"`
	CreatedAt   time.Time `form:"-" json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `form:"-" json:"updated_at" db:"updated_at"`
	//calculated fields
	Author       User     `form:"-" json:"author" db:"author"`
	Tags         []string `form:"tags" json:"tags" db:"-"` //can't make gin Bind form field to []Tag, so use []string instead
	CommentCount int64    `form:"-" json:"comment_count" db:"comment_count"`
}

//Insert saves Post as well as associated tags (creating them if needed) into db. Obsolete associations are removed
func (post *Post) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	err = db.QueryRow("INSERT INTO posts(name, description, published, user_id, created_at, updated_at) VALUES($1,$2,$3,$4,$5,$5) RETURNING id", post.Name, post.Description, post.Published, post.UserID, time.Now()).Scan(&post.ID)
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

//Update saves Post and associated tags changes into db
func (post *Post) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE posts SET name=$2, description=$3, published=$4, updated_at=$5 WHERE id=$1", post.ID, post.Name, post.Description, post.Published, time.Now())
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

//UpdateTags inserts new (non existent) post tags and updates associations
func (post *Post) UpdateTags(tx *sqlx.Tx) error {
	neu := make(map[string]bool)
	old := make(map[string]bool)
	for i := range post.Tags {
		neu[post.Tags[i]] = true
	}
	var exist []string
	err := db.Select(&exist, "SELECT tag_name FROM poststags WHERE post_id=$1", post.ID)
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
		_, err := tx.Exec("INSERT INTO poststags(post_id, tag_name) VALUES($1,$2)", post.ID, name)
		if err != nil {
			return err
		}
	}

	for name := range old {
		//remove association
		_, err := tx.Exec("DELETE FROM poststags WHERE post_id=$1 AND tag_name=$2", post.ID, name)
		if err != nil {
			return err
		}
	}
	return nil
}

//Delete removes Post from db. Existing postgresql contstraints remove tag associations on post delete.
func (post *Post) Delete() error {
	_, err := db.Exec("DELETE FROM posts WHERE id=$1", post.ID)
	return err
}

//Excerpt returns post excerpt by removing html tags first and truncating to 300 symbols
func (post *Post) Excerpt() template.HTML {
	//you can sanitize, cut it down, add images, etc
	policy := bluemonday.StrictPolicy() //remove all html tags
	sanitized := policy.Sanitize(string(blackfriday.MarkdownCommon([]byte(post.Description))))
	excerpt := template.HTML(truncate(sanitized, 300) + "...")
	return excerpt
}

//GetPost returns Post by its ID. Also initializes post author and tags fields.
func GetPost(id interface{}) (*Post, error) {
	post := &Post{}
	err := db.Get(post, "SELECT * FROM posts WHERE id=$1", id)
	if err != nil {
		return post, err
	}
	err = db.Get(&post.Author, "SELECT id,name FROM users WHERE id=$1", post.UserID)
	if err != nil {
		return post, err
	}
	err = db.Select(&post.Tags, "SELECT name FROM tags WHERE EXISTS (SELECT null FROM poststags WHERE post_id=$1 AND tag_name=tags.name)", id)
	return post, err
}

//GetPosts returns a slice of posts, order by descending id
func GetPosts() ([]Post, error) {
	var list []Post
	err := db.Select(&list, "SELECT * FROM posts ORDER BY posts.id DESC")
	return list, err
}

//GetPublishedPosts returns a slice published of posts with their associations
func GetPublishedPosts() ([]Post, error) {
	var list []Post
	err := db.Select(&list, "SELECT * FROM posts WHERE published=$1 ORDER BY posts.id DESC", true)
	if err != nil {
		return list, err
	}
	if err := fillPostsAssociations(list); err != nil {
		return list, err
	}
	return list, err
}

//GetRecentPosts returns a slice of published posts
func GetRecentPosts() ([]Post, error) {
	var list []Post
	err := db.Select(&list, "SELECT id, name FROM posts WHERE published=$1 ORDER BY id DESC LIMIT 7", true)
	return list, err
}

//GetPostMonths returns a slice of distinct months extracted from posts creation dates
func GetPostMonths() ([]Post, error) {
	var list []Post
	err := db.Select(&list, "SELECT DISTINCT date_trunc('month', created_at) as created_at FROM posts WHERE published=$1 ORDER BY created_at DESC", true)
	return list, err
}

//GetPostsByArchive returns a slice of published posts, given creation year and month
func GetPostsByArchive(year, month int) ([]Post, error) {
	var list []Post
	err := db.Select(&list, "SELECT * FROM posts WHERE published=$1 AND date_part('year', created_at)=$2 AND date_part('month', created_at)=$3 ORDER BY created_at DESC", true, year, month)
	if err != nil {
		return list, err
	}
	if err := fillPostsAssociations(list); err != nil {
		return list, err
	}
	return list, err
}

//GetPostsByTag returns a slice of published posts associated with tag name
func GetPostsByTag(name string) ([]Post, error) {
	var list []Post
	err := db.Select(&list, "SELECT * FROM posts WHERE published=$1 AND EXISTS (SELECT null FROM poststags WHERE poststags.post_id=posts.id AND poststags.tag_name=$2) ORDER BY created_at DESC", true, name)
	if err != nil {
		return list, err
	}
	if err := fillPostsAssociations(list); err != nil {
		return list, err
	}
	return list, nil
}

//fillPostsAssociations initialises post associations, given post slice
func fillPostsAssociations(list []Post) error {
	for i := range list {
		err := db.Get(&list[i].Author, "SELECT id,name FROM users WHERE id=$1", list[i].UserID)
		if err != nil {
			return err
		}
		err = db.Select(&list[i].Tags, "SELECT name FROM tags WHERE EXISTS (SELECT null FROM poststags WHERE post_id=$1 AND tag_name=tags.name)", list[i].ID)
		if err != nil {
			return err
		}
	}
	return nil
}
