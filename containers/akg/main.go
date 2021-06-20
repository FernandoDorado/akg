package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	k := &k8s{}
	k.InCluster = (os.Getenv("IN_CLUSTER") == "true")
	k.CloudProvider = os.Getenv("CLOUD_PROVIDER")
	if ok, _ := k.configure(); ok {
		k.connect()
	}

	// http
	e := gin.Default()
	e.LoadHTMLGlob("templates/*")
	e.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	e.GET("/api/v1/k8s/reconfigure", func(c *gin.Context) {
		ok, err := k.configure()
		if !ok {
			msg := fmt.Sprintf("failed to reconfigure kubernetes client: %s", err)
			log.Print(msg)
			c.JSON(500, gin.H{
				"error": msg,
			})
			return
		}
		msg := "reconfigured kubernetes client"
		log.Print(msg)
		c.JSON(200, gin.H{})
	})

	e.GET("/api/v1/k8s/reconnect", func(c *gin.Context) {
		ok, err := k.connect()
		if !ok {
			msg := fmt.Sprintf("failed to reconnect kubernetes client: %s", err)
			log.Print(msg)
			c.JSON(500, gin.H{
				"error": msg,
			})
			return
		}
		msg := "reconnected kubernetes client"
		log.Print(msg)
		c.JSON(200, gin.H{})
	})

	e.GET("/api/v1/k8s/apps", func(c *gin.Context) {
		if !k.Ready {
			c.JSON(500, gin.H{
				"error": "kubernetes client is not ready",
			})
			return
		}

		apps, err := k.apps()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"apps": apps,
		})
	})
	e.GET("/health/live", func(c *gin.Context) {})
	e.Run()
}
