package main

import (
    "github.com/gin-gonic/gin"
)

type Todo struct {
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

var todos = []Todo{}

func main() {
    router := gin.Default()
    router.LoadHTMLGlob("templates/index.html")

    router.GET("/", func(ctx *gin.Context) {
        ctx.HTML(200, "index.html", gin.H{
            "todos": todos,
        })
    })

    router.Run()
}