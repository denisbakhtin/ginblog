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
	"github.com/denisbakhtin/ginbasic/controllers"
	"github.com/denisbakhtin/ginbasic/controllers/admin"
	"github.com/denisbakhtin/ginbasic/helpers"
	"github.com/denisbakhtin/ginbasic/models"
	"github.com/denisbakhtin/ginbasic/system"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
)

func main() {
	migration := flag.String("migrate", "", "Run DB migrations: up, down, redo, new [MIGRATION_NAME] and then os.Exit(0)")
	flag.Parse()

	setLogger()
	loadConfig()
	connectToDB()
	runMigrations(migration)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	setTemplate(router)
	setSessions(router)

	router.StaticFS("/uploads", http.Dir(system.GetConfig().Uploads))
	router.StaticFS("/public", rice.MustFindBox("public").HTTPBox()) //<3 rice

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

	authorized := router.Group("/admin")
	authorized.Use(AuthRequired())
	authorized.GET("/", admin.AdminGet)

	authorized.POST("/upload", admin.UploadPost) //image upload

	authorized.GET("/users", admin.UserIndex)
	authorized.GET("/new_user", admin.UserNew)
	authorized.POST("/new_user", admin.UserCreate)
	authorized.GET("/users/:id/edit", admin.UserEdit)
	authorized.POST("/users/:id/edit", admin.UserUpdate)
	authorized.POST("/users/:id/delete", admin.UserDelete)

	authorized.GET("/pages", admin.PageIndex)
	authorized.GET("/new_page", admin.PageNew)
	authorized.POST("/new_page", admin.PageCreate)
	authorized.GET("/pages/:id/edit", admin.PageEdit)
	authorized.POST("/pages/:id/edit", admin.PageUpdate)
	authorized.POST("/pages/:id/delete", admin.PageDelete)
	/*
	   router.GET("/someGet", getting)
	   router.POST("/somePost", posting)
	   router.PUT("/somePut", putting)
	   router.DELETE("/someDelete", deleting)
	   router.PATCH("/somePatch", patching)
	   router.HEAD("/someHead", head)
	   router.OPTIONS("/someOptions", options)
	*/

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
	models.SetDB(system.GetConfig())
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
		"isActive": helpers.IsActive,
		"dateTime": helpers.DateTime,
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
	store := sessions.NewCookieStore([]byte(config.SessionSecret))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("gin-session", store))
	//https://github.com/utrack/gin-csrf
	router.Use(csrf.Middleware(csrf.Options{
		Secret: config.CsrfSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))
}

//+++++++++++++ middlewares +++++++++++++++++++++++
//SharedData fills in common data, such as user info, etc...
func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if uId := session.Get("UserId"); uId != nil {
			user, _ := models.GetUser(uId)
			if user.Id != 0 {
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
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("Access forbidden"))
		}
	}
}
