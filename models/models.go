package models

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/fiam/gounidecode/unidecode"
	"gopkg.in/loremipsum.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Model is a tuned version of gorm.Model
type Model struct {
	ID        uint64     `form:"id" gorm:"primary_key"`
	CreatedAt time.Time  `binding:"-" form:"-"`
	UpdatedAt time.Time  `binding:"-" form:"-"`
	DeletedAt *time.Time `binding:"-" form:"-"`
}

var db *gorm.DB

// SetDB establishes connection to database and saves its handler into db *sqlx.DB
func SetDB(connection string) {
	var err error
	db, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

// GetDB returns database handler
func GetDB() *gorm.DB {
	return db
}

// AutoMigrate runs gorm auto migration
func AutoMigrate() {
	if err := db.AutoMigrate(&User{}, &Tag{}, &Page{}, &Post{}, &Comment{}); err != nil {
		panic(err)
	}
}

// Seed initializes empty database with initial data
func SeedDB() {
	var users []User
	db.Find(&users)
	doSeed := len(users) == 0
	if doSeed {
		lorem := loremipsum.New()
		user := User{Email: "admin@email.com", Name: "admin", Password: "admin"}
		if err := db.Create(&user).Error; err != nil {
			slog.Error(err.Error())
		}
		post := Post{
			Title:     "Welcome to ginblog",
			Content:   fmt.Sprintf(`<p>Sit back and relax.</p><p><img src="/public/images/welcome_image.jpg" /></p><p>%s</p><p>%s</p>`, lorem.Paragraph(), lorem.Paragraph()),
			UserID:    user.ID,
			Published: true,
			Comments:  []Comment{{UserName: "Admin", Content: fmt.Sprintf("Keep it up! %s.", lorem.Sentence()), Published: true, UserID: user.ID}},
			Tags:      []Tag{{Title: "intro"}, {Title: "seeding DB"}},
		}
		if err := db.Create(&post).Error; err != nil {
			slog.Error(err.Error())
		}
		page := Page{
			Title:     "Contacts",
			Content:   fmt.Sprintf("<p>South Pole, Warm Street 1.</p><p>%s</p><p>%s</p>", lorem.Paragraph(), lorem.Paragraph()),
			Published: true,
		}
		if err := db.Create(&page).Error; err != nil {
			slog.Error(err.Error())
		}
	}
}

// truncate truncates string to n runes
func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}

// createSlug makes a friendly url slug out of a string
func createSlug(s string) string {
	s = strings.ToLower(unidecode.Unidecode(s))                    //transliterate if it is not in english
	s = regexp.MustCompile(`[^a-z0-9\s]+`).ReplaceAllString(s, "") //remove non alpha-numeric characters
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, "-")         //replace all spaces with dashes
	return s
}
