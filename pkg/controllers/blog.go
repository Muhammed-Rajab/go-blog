package controllers

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/Muhammed-Rajab/go-blog/pkg/db"
	"github.com/Muhammed-Rajab/go-blog/pkg/models"
	"github.com/Muhammed-Rajab/go-blog/pkg/utils"
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

	pageNo := ctx.Query("page")
	search := ctx.Query("search")

	page, err := strconv.ParseUint(pageNo, 10, 64)
	if err != nil {
		page = 1
	}

	posts, err := blogs.FindBlogs(bson.M{
		"title": bson.M{
			"$regex":   search,
			"$options": "i",
		},
		"published": true,
	}, int(page), 10)

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
	templates, _ := template.New("custom-blog").Funcs(utils.GetTemplateFuncsMap()).ParseGlob("templates/*")

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

func (BlogController) DashboardHandler(ctx *gin.Context) {

	obj := gin.H{}
	blogs := models.NewBlogs(db.GetMDB().BlogsCollection())

	pageNo := ctx.Query("page")
	search := ctx.Query("search")

	page, err := strconv.ParseUint(pageNo, 10, 64)
	if err != nil {
		page = 1
	}

	posts, err := blogs.FindBlogs(bson.M{
		"title": bson.M{
			"$regex":   search,
			"$options": "i",
		},
	}, int(page), 10)

	if err != nil {
		obj["errors"] = err
	} else {
		obj["posts"] = posts
	}
	ctx.HTML(http.StatusOK, "dashboard.html", obj)
}
