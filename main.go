package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/home")
	})

	router.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"title": "Xoom website",
		})
	})

	router.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base.tmpl", gin.H{
			"title": "testpage",
		})
	})

	router.Run(":8080")
}
