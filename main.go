package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"

	_ "github.com/go-sql-driver/mysql"
)

var (
	env     string
	headers gin.H = gin.H{}
)

func main() {

	// db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/db_name")
	// if err != nil {
	// 	fmt.Printf("Error connecting to db: %v\n", err)
	// }

	// rows, err := db.Query("SELECT * FROM about")
	// if err != nil {
	// 	fmt.Printf("some error -> %v\n", err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var output string

	// 	err := rows.Scan(&output)
	// 	if err != nil {
	// 		fmt.Printf("some error -> %v\n", err)
	// 	}

	// 	fmt.Printf("Value retrieved from db -> %v\n", output)

	// }

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

	router.GET("/interests", func(c *gin.Context) {
		c.HTML(http.StatusOK, "interests.tmpl", headers)
	})

	router.GET("/blog", func(c *gin.Context) {
		c.HTML(http.StatusOK, "blog.tmpl", headers)
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
