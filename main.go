package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	// Root router for the application
	root := gin.Default()

	// Routing starts here
	blog := root.Group("/blog")
	{
		blog.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "welcome to blog!",
			})
		})
	}

	root.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hello world!",
		})
	})

	root.Run(":8000")
}
