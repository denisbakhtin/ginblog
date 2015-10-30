package models

import (
	"html/template"
	"time"

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
	AuthorName null.String `form:"-" json:"author_name" db:"author_name"`
}

func (post *Post) Insert() error {
	err := db.QueryRow("INSERT INTO posts(name, description, published, user_id, timestamp) VALUES($1,$2,$3,$4,$5) RETURNING id", post.Name, post.Description, post.Published, post.UserId, time.Now()).Scan(&post.Id)
	return err
}

func (post *Post) Update() error {
	_, err := db.Exec("UPDATE posts SET name=$2, description=$3, published=$4 WHERE id=$1", post.Id, post.Name, post.Description, post.Published)
	return err
}

func (post *Post) Delete() error {
	_, err := db.Exec("DELETE FROM posts WHERE id=$1", post.Id)
	return err
}

func (post *Post) Excerpt() template.HTML {
	//you can sanitize, cut it down, add images, etc
	policy := bluemonday.StrictPolicy() //remove all html tags
	sanitized := policy.Sanitize(string(blackfriday.MarkdownCommon([]byte(post.Description))))
	excerpt := template.HTML(sanitized)
	return excerpt
}

func GetPost(id interface{}) (*Post, error) {
	post := &Post{}
	err := db.Get(post, "SELECT posts.*, users.name as author_name FROM posts LEFT OUTER JOIN users ON posts.user_id=users.id WHERE posts.id=$1", id)
	return post, err
}

func GetPosts() ([]Post, error) {
	var list []Post
	err := db.Select(&list, "SELECT posts.*, users.name as author_name FROM posts LEFT OUTER JOIN users ON posts.user_id=users.id ORDER BY posts.id DESC")
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

func GetPostsByYearMonth(year, month int) ([]Post, error) {
	var list []Post
	err := db.Select(&list, "SELECT * FROM posts WHERE published=$1 AND date_part('year', timestamp)=$2 AND date_part('month', timestamp)=$3 ORDER BY timestamp DESC", true, year, month)
	return list, err
}
