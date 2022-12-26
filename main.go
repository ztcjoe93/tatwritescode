package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"internal/database"
	"internal/posts"
	"internal/utilities"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

var (
	env         string
	headers     gin.H  = gin.H{}
	identityKey string = "id"
	db          *sql.DB
	blogposts   []*posts.Blogpost
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}

	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbHost := os.Getenv("MYSQL_HOST")
	env = os.Getenv("ENV")

	if env == "" {
		dbHost = "localhost"
		env = "dev"
	}

	db := database.OpenSqlConnection(dbUser, dbPassword, dbName, dbHost)
	blogposts = database.GetAllBlogposts(db)

	// For Cache-Control: "no-cache" when running on development so that
	// we don't get a mysterious 'why isn't my assets updating' situation
	if env == "dev" {
		headers["Cache-Control"] = "no-cache"
		fmt.Println("Cache-Control set to 'no-cache'")
	} else {
		headers = gin.H{}
	}

	router := gin.New()

	funcMap := template.FuncMap{
		"monthIntRepr": utilities.ConvertMonthToIntRepr,
		"renderAsHTML": utilities.RenderAsHTML,
	}
	router.SetFuncMap(funcMap)

	router.LoadHTMLGlob("templates/*")

	router.Static("/js", "./js/")
	router.Static("/assets", "./assets/")

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		SendCookie:     true,
		SecureCookie:   env == "prod", //non HTTPS dev environments
		CookieHTTPOnly: true,          // JS can't modify
		CookieName:     "token",       // default jwt
		TokenLookup:    "cookie:token",
		CookieSameSite: http.SameSiteDefaultMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode
		Realm:          "TATWRITESCODE.COM",
		Key:            []byte(os.Getenv("SIGNATURE_KEY")),
		Timeout:        time.Hour,
		MaxRefresh:     time.Hour,
		IdentityKey:    identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			hashedPassword := database.GetHashedPassword(db, username)
			fmt.Printf("checkpass = %v\n", utilities.CheckPasswordHash(password, hashedPassword))

			if utilities.CheckPasswordHash(password, hashedPassword) {
				return &User{
					UserName: username,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*User); ok {
				fmt.Printf("Data = %v, ok = %v\n", data, ok)
				return true
			}
			return false
		},
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			c.Redirect(http.StatusFound, "/admin/home")
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.Redirect(http.StatusFound, "/")
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		// TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/logout", authMiddleware.LogoutHandler)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.tmpl", headers)
	})

	router.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.tmpl", headers)
	})

	router.GET("/experience", func(c *gin.Context) {
		c.HTML(http.StatusOK, "experience.tmpl", headers)
	})

	router.GET("/education", func(c *gin.Context) {
		c.HTML(http.StatusOK, "education.tmpl", headers)
	})

	router.GET("/skills", func(c *gin.Context) {
		c.HTML(http.StatusOK, "skills.tmpl", headers)
	})

	router.GET("/projects", func(c *gin.Context) {
		c.HTML(http.StatusOK, "projects.tmpl", headers)
	})

	router.GET("/blog", func(c *gin.Context) {
		latestPosts := database.GetLatestPosts(db)

		c.HTML(http.StatusOK, "blog.tmpl", gin.H{
			"posts":   latestPosts,
			"links":   posts.GetNavigationLinks(blogposts),
			"baseUrl": getBaseURL(c),
		})
	})

	router.GET("/blog/:year/:month", func(c *gin.Context) {
		year := c.Param("year")
		month := c.Param("month")

		specificPosts, _ := database.GetPostsFromMonth(db, year, month)

		c.HTML(http.StatusOK, "blog.tmpl", gin.H{
			"posts":   specificPosts,
			"links":   posts.GetNavigationLinks(blogposts),
			"baseUrl": getBaseURL(c, true),
		})
	})

	router.GET("/resume", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "assets/resume.pdf")
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", headers)
	})

	router.GET("/logout", func(c *gin.Context) {
		authMiddleware.LogoutHandler(c)
		c.Redirect(http.StatusFound, "/")
	})
	auth := router.Group("/admin")

	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	auth.GET("/home", AuthHandler)
	auth.POST("/home", func(c *gin.Context) {
		title := c.PostForm("post_title")
		post := c.PostForm("post_content")

		database.InsertPost(db, time.Now().UTC().Format("2006-01-02 03:04:05"), title, post)
		c.Redirect(http.StatusFound, "/admin/home")
	})

	router.Run(":8080")
}

/*
	Modify headers with a map of arguments and return a new gin.H using the base headers
*/
func modH(args map[string]string) gin.H {
	header := headers

	for key, value := range args {
		header[key] = value
	}

	return header
}

func getBaseURL(c *gin.Context, additionalParams ...bool) string {
	scheme := "http"

	if c.Request.TLS != nil {
		scheme = "https"
	}

	urlPath := c.Request.URL.Path
	if len(additionalParams) > 0 {
		arguments := strings.Split(c.Request.URL.Path, "/")
		fmt.Printf("%v\n", arguments)
		urlPath = "/" + arguments[1]
	}

	return scheme + "://" + c.Request.Host + urlPath
}

func AuthHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	// user, _ := c.Get(identityKey)
	fmt.Printf("%v\n", claims)
	c.HTML(http.StatusOK, "admin.tmpl", headers)
}

type User struct {
	UserName string
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
