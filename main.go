package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Muhammed-Rajab/go-blog/pkg/db"
	"github.com/Muhammed-Rajab/go-blog/pkg/routers"
	"github.com/Muhammed-Rajab/go-blog/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("./.env"); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	// Making sure that there's a BLOG_DASHBOARD_KEY on env
	if os.Getenv("BLOG_DASHBOARD_KEY") == "" {
		log.Fatalf("BLOG_DASHBOARD_KEY missing in env variable")
	}

	db.Init("mongodb://localhost:27017")
	db.GetMDB().Connect()

	engine := gin.Default()
	engine.SetFuncMap(utils.GetTemplateFuncsMap())
	engine.LoadHTMLGlob("templates/*.html")
	engine.Static("/public", "./public")
	engine.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "notfound.html", nil)
	})

	root := engine.Group("/")
	routers.BlogRouter(root)

	engine.Run(":8000")
}
