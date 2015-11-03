//go:generate rice embed-go
package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/GeertJohan/go.rice"
	"github.com/Sirupsen/logrus"
	"github.com/claudiu/gocron"
	"github.com/denisbakhtin/ginblog/controllers"
	"github.com/denisbakhtin/ginblog/helpers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/denisbakhtin/ginblog/system"
	//"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	//"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	//"github.com/utrack/gin-csrf"
)

const (
	sessionName = "session"
)

func main() {
	migration := flag.String("migrate", "", "Run DB migrations: up, down, redo, new [MIGRATION_NAME] and then os.Exit(0)")
	flag.Parse()

	setLogger()
	loadConfig()
	connectToDB()
	runMigrations(migration)

	//Periodic tasks
	gocron.Every(1).Day().Do(system.CreateXMLSitemap)
	gocron.Start()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	setTemplate(router) //initialize templates
	setSessions(router) //initialize session storage & use sessiom/csrf middlewares

	router.StaticFS("/public", http.Dir(system.PublicPath())) //better use nginx to serve assets (Cache-Control, Etag, fast gzip, etc)
	router.Use(SharedData())

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

	router.GET("/pages/:id", controllers.PageGet)
	router.GET("/posts/:id", controllers.PostGet)
	router.GET("/tags/:name", controllers.TagGet)
	router.GET("/archives/:year/:month", controllers.ArchiveGet)
	router.GET("/rss", controllers.RssGet)

	authorized := router.Group("/admin")
	authorized.Use(AuthRequired())
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

		authorized.GET("/tags", controllers.TagIndex)
		authorized.GET("/new_tag", controllers.TagNew)
		authorized.POST("/new_tag", controllers.TagCreate)
		authorized.POST("/tags/:name/delete", controllers.TagDelete)
	}

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}

//setLogger initializes logrus logger with some defaults
func setLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stderr)
	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

//setConfig loads config.json from rice box "config"
func loadConfig() {
	box := rice.MustFindBox("config")
	system.LoadConfig(box.MustBytes("config.json"))
}

//connectToDB initializes *sqlx.DB handler
func connectToDB() {
	config := system.GetConfig()
	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.User, config.Database.Password, config.Database.Name)
	models.SetDB(connection)
}

//runMigrations applies database migrations if "migrate" flat is set
func runMigrations(command *string) {
	if len(*command) > 0 {
		//Read https://github.com/rubenv/sql-migrate for more info about migrations
		box := rice.MustFindBox("migrations")
		system.RunMigrations(box, models.GetDB(), command)
		os.Exit(0)
	}
}

//setTemplate loads templates from rice box "views"
func setTemplate(router *gin.Engine) {
	box := rice.MustFindBox("views")
	tmpl := template.New("").Funcs(template.FuncMap{
		"isActive":      helpers.IsActive,
		"stringInSlice": helpers.StringInSlice,
		"dateTime":      helpers.DateTime,
		"recentPosts":   helpers.RecentPosts,
		"tags":          helpers.Tags,
		"archives":      helpers.Archives,
	})

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			var err error
			tmpl, err = tmpl.Parse(box.MustString(path))
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := box.Walk("", fn)
	if err != nil {
		panic(err)
	}
	router.SetHTMLTemplate(tmpl)
}

//setSessions initializes sessions & csrf middlewares
func setSessions(router *gin.Engine) {
	config := system.GetConfig()
	//https://github.com/gin-gonic/contrib/tree/master/sessions
	//TODO: switch to pure gorilla/sessions & gorilla/csrf
	store := sessions.NewFilesystemStore("", []byte(config.SessionSecret))
	store.Options = &sessions.Options{HttpOnly: true, MaxAge: 7 * 86400} //Also set Secure: true if using SSL, you should though
	router.Use(SessionMiddleware("session", store))
	//router.Use(gin.WrapH(csrf.Protect([]byte(config.SessionSecret))(router)))
	//https://github.com/utrack/gin-csrf
	/*
		router.Use(csrf.Middleware(csrf.Options{
			Secret: config.SessionSecret,
			ErrorFunc: func(c *gin.Context) {
				c.String(400, "CSRF token mismatch")
				c.Abort()
			},
		}))
	*/
}

//+++++++++++++ middlewares +++++++++++++++++++++++

//SessionMiddleware stores session handler in gin.Context
func SessionMiddleware(name string, store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer context.Clear(c.Request)
		session, err := store.Get(c.Request, name)
		if err != nil {
			logrus.Error(err)
			c.AbortWithError(500, err)
			return
		}
		c.Set("Session", session)
		c.Next()
	}
}

//SharedData fills in common data, such as user info, etc...
func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := helpers.GetSession(c)
		if uID := session.Values["UserID"]; uID != nil {
			user, _ := models.GetUser(uID)
			if user.ID != 0 {
				c.Set("User", user)
			}
		}
		if system.GetConfig().SignupEnabled {
			c.Set("SignupEnabled", true)
		}
		c.Next()
	}
}

//AuthRequired grants access to authenticated users, requires SharedData middleware
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("User"); user != nil {
			c.Next()
		} else {
			logrus.Warnf("User not authorized to visit %s", c.Request.RequestURI)
			c.HTML(http.StatusForbidden, "errors/403", nil)
			c.Abort()
		}
	}
}
