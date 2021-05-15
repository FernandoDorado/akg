package webserver

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type WebServer struct {
    Engine *gin.Engine
}

func (ws *WebServer) init() {
    ws.Engine = gin.Default()
    ws.Engine.LoadHTMLGlob("templates/*")
    ws.Engine.GET("/health/live", func(c *gin.Context) {})
    ws.Engine.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H{})
    })
}

func (ws *WebServer) run() {
    ws.Engine.Run()
}

func (ws *WebServer) Start() {
    ws.init()
    ws.run()
}
