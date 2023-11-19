package main

import (
	"log"
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

	root := engine.Group("/")
	routers.BlogRouter(root)

	engine.Run(":8000")

	// blogs := models.NewBlogs(db.GetMDB().BlogsCollection())
	// post := models.BlogModel{
	// 	Title: "Far, but close.",
	// 	Desc:  "How can someone be so far away from you physically, yet so close to you spiritually?",
	// 	Content: `<p style="text-align: center;">
	// Amidst of all the chaos and fun,<br>
	// I could only think of you.<br>
	// It never fails to intrigue me,<br>
	// How do you live inside my head rent free?
	// <br>
	// <br>
	// We are at least a century miles away,<br>
	// Yet I could smell you here,<br>
	// feel your presence; feel your forehead,<br>
	// getting warm as I leave a kiss mark on it.
	// <br>
	// <br>
	// Four days just seem trivial, but are they?
	// I guess not.<br>
	// </p>`,
	// 	Tags: []string{"love", "distance", "relationship"},
	// }

	// if err := blogs.AddBlog(post); err != nil {
	// 	log.Fatal(err)
	// }
}
