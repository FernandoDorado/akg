package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()
	e.LoadHTMLGlob("templates/*")
	e.GET("/", func(c *gin.Context) {
		if strings.Contains(c.Request.Header.Get("User-Agent"), "curl") {
			c.JSON(http.StatusOK, gin.H{
				"Name":  "Adam K Gray",
				"About": "Solutions Architecture, Data Science",
			})
			return
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	e.GET("/health/live", func(c *gin.Context) {})
	e.Run()
}
