package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"internal/database"
	"internal/posts"
	"internal/utilities"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"

	_ "github.com/go-sql-driver/mysql"
)

var (
	env       string
	headers   gin.H = gin.H{}
	db        *sql.DB
	blogposts []*posts.Blogpost
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbHost := os.Getenv("MYSQL_HOST")

	if dbHost == "" {
		dbHost = "localhost"
	}

	db := database.OpenSqlConnection(dbUser, dbPassword, dbName, dbHost)
	blogposts = database.GetAllBlogposts(db)

	config := make(map[interface{}]interface{})
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Error reading yaml file: %v\n", err)
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Printf("Unable to unmarshal config.yaml: %v\n", err)
	}

	if _, ok := config["env"]; ok {
		env = config["env"].(string)
	} else {
		env = "dev"
	}

	// For Cache-Control: "no-cache" when running on development so that
	// we don't get a mysterious 'why isn't my assets updating' situation
	if env == "dev" {
		headers["Cache-Control"] = "no-cache"
		fmt.Println("Cache-Control set to 'no-cache'")
	} else {
		headers = gin.H{}
	}

	fmt.Printf("Currently in %v environment\n", env)

	router := gin.Default()

	funcMap := template.FuncMap{
		"monthIntRepr": utilities.ConvertMonthToIntRepr,
	}
	router.SetFuncMap(funcMap)

	router.LoadHTMLGlob("templates/*")

	router.Static("/js", "./js/")
	router.Static("/assets", "./assets/")

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
