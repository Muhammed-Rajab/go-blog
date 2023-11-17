package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	// Root router for the application
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	// Routing starts here
	blog := router.Group("/blog")
	{
		blog.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "welcome to blog!",
			})
		})
	}

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"message": "hello!",
		})
	})

	router.Run(":8000")
}
