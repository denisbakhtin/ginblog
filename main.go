package main

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/claudiu/gocron"
	"github.com/denisbakhtin/ginblog/controllers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/denisbakhtin/ginblog/system"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
)

func main() {
	initLogger()
	system.LoadConfig()
	models.SetDB(system.GetConnectionString())
	models.AutoMigrate()
	system.LoadTemplates()

	//Periodic tasks
	gocron.Every(1).Day().Do(system.CreateXMLSitemap)
	gocron.Start()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.SetHTMLTemplate(system.GetTemplates())

	//setup sessions
	config := system.GetConfig()
	store := memstore.NewStore([]byte(config.SessionSecret))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("gin-session", store))

	//setup csrf protection
	router.Use(csrf.Middleware(csrf.Options{
		Secret: config.SessionSecret,
		ErrorFunc: func(c *gin.Context) {
			logrus.Error("CSRF token mismatch")
			controllers.ShowErrorPage(c, 400, fmt.Errorf("CSRF token mismatch"))
			c.Abort()
		},
	}))

	router.StaticFS("/public", http.Dir(system.PublicPath())) //better use nginx to serve assets (Cache-Control, Etag, fast gzip, etc)
	router.Use(controllers.ContextData())

	router.GET("/", controllers.HomeGet)
	router.NoRoute(controllers.NotFound)
	router.NoMethod(controllers.MethodNotAllowed)

	if system.GetConfig().SignupEnabled {
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
	router.GET("/tags/:title", controllers.TagGet)
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
		authorized.POST("/tags/:title/delete", controllers.TagDelete)
	}

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}

//initLogger initializes logrus logger with some defaults
func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	//logrus.SetOutput(os.Stderr)
	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
