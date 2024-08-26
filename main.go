package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/claudiu/gocron"
	"github.com/denisbakhtin/ginblog/config"
	"github.com/denisbakhtin/ginblog/controllers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func main() {
	initLogger()
	config.LoadConfig()
	models.SetDB(config.GetConnectionString())
	models.AutoMigrate()
	models.SeedDB() // seed database with demo data if empty
	controllers.LoadTemplates()

	//Periodic tasks
	gocron.Every(1).Day().Do(controllers.CreateXMLSitemap)
	gocron.Start()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.SetHTMLTemplate(controllers.GetTemplates())

	//setup sessions
	store := memstore.NewStore([]byte(config.GetConfig().SessionSecret))
	store.Options(sessions.Options{Path: "/", HttpOnly: true, MaxAge: 7 * 86400}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("ginblog-session", store))

	//setup csrf protection
	router.Use(csrf.Middleware(csrf.Options{
		Secret: config.GetConfig().SessionSecret,
		ErrorFunc: func(c *gin.Context) {
			slog.Error("CSRF token mismatch")
			controllers.ShowErrorPage(c, 400, fmt.Errorf("CSRF token mismatch"))
			c.Abort()
		},
	}))

	router.StaticFS("/public", http.Dir(config.PublicPath())) //better use nginx to serve assets (Cache-Control, Etag, fast gzip, etc)
	router.Use(controllers.ContextData())

	router.GET("/", controllers.HomeGet)
	router.NoRoute(controllers.NotFound)
	router.NoMethod(controllers.MethodNotAllowed)

	if config.GetConfig().SignupEnabled {
		router.GET("/signup", controllers.SignUpGet)
		router.POST("/signup", controllers.SignUpPost)
	}
	router.GET("/signin", controllers.SignInGet)
	router.POST("/signin", controllers.SignInPost)
	router.GET("/logout", controllers.LogoutGet)

	router.GET("/oauthgooglelogin", controllers.OauthGoogleLogin)
	router.GET("/oauthcallback", controllers.OauthCallback)
	router.POST("/new_comment", controllers.CommentPublicCreate)

	router.GET("/pages/:id", controllers.PageGet)
	router.GET("/posts/:id", controllers.PostGet)
	router.GET("/tags/:slug", controllers.TagGet)
	router.GET("/archives/:year/:month", controllers.ArchiveGet)
	router.GET("/rss", controllers.RssGet)

	authorized := router.Group("/admin")
	authorized.Use(controllers.AuthRequired())
	{
		authorized.GET("/", controllers.AdminGet)

		authorized.POST("/upload", controllers.UploadPost) //image upload

		authorized.GET("/users", controllers.UserIndex)
		authorized.GET("/new_user", controllers.UserNew)
		authorized.POST("/new_user", controllers.UserCreate)
		authorized.GET("/users/:id/edit", controllers.UserEdit)
		authorized.POST("/users/:id/edit", controllers.UserUpdate)
		authorized.POST("/users/:id/delete", controllers.UserDelete)

		authorized.GET("/pages", controllers.PageIndex)
		authorized.GET("/new_page", controllers.PageNew)
		authorized.POST("/new_page", controllers.PageCreate)
		authorized.GET("/pages/:id/edit", controllers.PageEdit)
		authorized.POST("/pages/:id/edit", controllers.PageUpdate)
		authorized.POST("/pages/:id/delete", controllers.PageDelete)

		authorized.GET("/posts", controllers.PostIndex)
		authorized.GET("/new_post", controllers.PostNew)
		authorized.POST("/new_post", controllers.PostCreate)
		authorized.GET("/posts/:id/edit", controllers.PostEdit)
		authorized.POST("/posts/:id/edit", controllers.PostUpdate)
		authorized.POST("/posts/:id/delete", controllers.PostDelete)

		authorized.GET("/comments", controllers.CommentIndex)
		authorized.GET("/new_comment", controllers.CommentNew)
		authorized.POST("/new_comment", controllers.CommentCreate)
		authorized.GET("/comments/:id/edit", controllers.CommentEdit)
		authorized.POST("/comments/:id/edit", controllers.CommentUpdate)
		authorized.POST("/comments/:id/delete", controllers.CommentDelete)

		authorized.GET("/tags", controllers.TagIndex)
		authorized.GET("/new_tag", controllers.TagNew)
		authorized.POST("/new_tag", controllers.TagCreate)
		authorized.POST("/tags/:slug/delete", controllers.TagDelete)
	}

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}

// initLogger initializes slog logger with some defaults
func initLogger() {
	var handler slog.Handler
	if gin.Mode() == gin.ReleaseMode {
		opts := &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
