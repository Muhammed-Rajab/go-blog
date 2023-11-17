package main

import (
	"text/template"

	"github.com/Muhammed-Rajab/go-blog/pkg/db"
	"github.com/Muhammed-Rajab/go-blog/pkg/routers"
	"github.com/gin-gonic/gin"
)

func main() {

	db.Init("mongodb://localhost:27017")
	db.GetMDB().Connect()

	engine := gin.Default()
	engine.SetFuncMap(template.FuncMap{})
	engine.LoadHTMLGlob("templates/*.html")

	root := engine.Group("/")
	routers.BlogRouter(root)

	engine.Run(":8000")
}
