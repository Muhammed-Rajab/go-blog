package main

import (
	"github.com/Muhammed-Rajab/go-blog/pkg/db"
	"github.com/Muhammed-Rajab/go-blog/pkg/routers"
	"github.com/Muhammed-Rajab/go-blog/pkg/utils"
	"github.com/gin-gonic/gin"
)

func main() {

	db.Init("mongodb://localhost:27017")
	db.GetMDB().Connect()

	engine := gin.Default()
	engine.SetFuncMap(utils.GetTemplateFuncsMap())
	engine.LoadHTMLGlob("templates/*.html")
	engine.Static("/public", "./public")

	root := engine.Group("/")
	routers.BlogRouter(root)

	engine.Run(":8000")

	// 	blogs := models.NewBlogs(db.GetMDB().BlogsCollection())
	// 	post := models.BlogModel{
	// 		Title: "Far, but close.",
	// 		Content: `Amidst of all the chaos and fun,
	// I could only think of you.
	// It never fails to intrigue me,
	// How do you live inside my head rent free?

	// We are at least a century miles away,
	// Yet I could smell you here,
	// feel your presence; feel your forehead,
	// getting warm as I leave a kiss mark on it.

	// Four days just seem trivial, but are they?
	// I guess not.
	// `,
	// 		Tags: []string{"love", "hate"},
	// 	}

	//	if err := blogs.AddBlog(post); err != nil {
	//		log.Fatal(err)
	//	}
}
