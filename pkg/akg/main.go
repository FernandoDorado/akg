package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// http
	e := gin.Default()
	e.LoadHTMLGlob("templates/*")
	e.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	e.GET("/health/live", func(c *gin.Context) {})
	e.Run()
}
