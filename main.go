package main

import (
	"github.com/Muhammed-Rajab/go-blog/pkg/routers"
	"github.com/gin-gonic/gin"
)

func main() {

	engine := gin.Default()
	engine.LoadHTMLGlob("templates/*.html")

	root := engine.Group("/")
	routers.BlogRouter(root)

	engine.Run(":8000")
}
