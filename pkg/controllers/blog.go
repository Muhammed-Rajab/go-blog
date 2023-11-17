package controllers

import (
	"net/http"
	"text/template"

	"github.com/Muhammed-Rajab/go-blog/pkg/db"
	"github.com/Muhammed-Rajab/go-blog/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"go.mongodb.org/mongo-driver/bson"
)

type BlogController struct{}

func NewBlogController() *BlogController {
	return &BlogController{}
}

func (BlogController) HomeHandler(ctx *gin.Context) {

	obj := gin.H{}
	blogs := models.NewBlogs(db.GetMDB().BlogsCollection())

	posts, err := blogs.FindBlogs(bson.M{}, 1, 10)

	if err != nil {
		obj["errors"] = err
	} else {
		obj["posts"] = posts
	}
	ctx.HTML(http.StatusOK, "index.html", obj)
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func (BlogController) BlogHandler(ctx *gin.Context) {

	obj := gin.H{
		"author": "Rajab",
	}
	slug := ctx.Param("slug")
	blogs := models.NewBlogs(db.GetMDB().BlogsCollection())

	post, err := blogs.FindBlogBySlug(slug)
	templates := template.Must(template.ParseFiles("templates/blog.html"))

	if err != nil {
		obj["errors"] = err
	} else {
		post.Content = string(mdToHTML([]byte(post.Content)))
		obj["post"] = post
	}

	if err := templates.ExecuteTemplate(ctx.Writer, "blog.html", obj); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}

}
