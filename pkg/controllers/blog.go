package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/Muhammed-Rajab/go-blog/pkg/db"
	"github.com/Muhammed-Rajab/go-blog/pkg/models"
	"github.com/Muhammed-Rajab/go-blog/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (BlogController) TogglePublishBlog(ctx *gin.Context) {
	var obj gin.H

	id := ctx.Param("id")
	objectid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		obj = gin.H{
			"status": "failed",
			"error":  "invalid object id",
		}
		ctx.JSON(http.StatusBadRequest, obj)
		return
	}

	blogs := models.NewBlogs(db.GetMDB().BlogsCollection())

	if err := blogs.PublishDraftBlogByID(objectid.Hex()); err != nil {
		obj = gin.H{
			"status": "failed",
			"error":  "failed to toggle publish for blog:" + err.Error(),
		}
		ctx.JSON(http.StatusBadRequest, obj)
		return
	}

	obj = gin.H{
		"status":  "success",
		"message": "successfully toggle publish",
	}
	ctx.JSON(http.StatusOK, obj)
}

func (BlogController) DeleteBlog(ctx *gin.Context) {
	var obj gin.H

	id := ctx.Param("id")
	objectid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		obj = gin.H{
			"status": "failed",
			"error":  "invalid object id",
		}
		ctx.JSON(http.StatusBadRequest, obj)
		return
	}

	blogs := models.NewBlogs(db.GetMDB().BlogsCollection())

	if err := blogs.DeleteBlogByID(objectid.Hex()); err != nil {
		obj = gin.H{
			"status": "failed",
			"error":  "failed to delete blog:" + err.Error(),
		}
		ctx.JSON(http.StatusBadRequest, obj)
		return
	}

	obj = gin.H{
		"status":  "success",
		"message": "successfully deleted post",
	}
	ctx.JSON(http.StatusOK, obj)
}

func (BlogController) EditBlog(ctx *gin.Context) {

}

func (BlogController) AddBlogHandler(ctx *gin.Context) {
	var obj gin.H
	ctx.HTML(http.StatusOK, "add_blog.html", obj)
}

func (BlogController) AddBlog(ctx *gin.Context) {
	var obj gin.H = gin.H{}
	var form models.BlogForm

	if err := ctx.ShouldBind(&form); err != nil {
		obj["error"] = "Shit took a turn for the worst: " + err.Error()
		ctx.JSON(http.StatusBadRequest, obj)
		return
	}

	blogs := models.NewBlogs(db.GetMDB().BlogsCollection())

	// Create BlogModel from form
	publish := false
	if form.Publish == "on" {
		publish = true
	}

	blog := models.BlogModel{
		Title:     form.Title,
		Desc:      form.Desc,
		Content:   form.Content,
		Tags:      tagsFromString(form.Tags),
		Published: publish,
	}

	// Redirect to the created blog if everything went well
	ctx.Redirect(http.StatusSeeOther, "/blog")
}

func tagsFromString(stringTag string) []string {
	tags := strings.Split(stringTag, ",")
	for i, tag := range tags {
		tags[i] = strings.ToLower(strings.TrimSpace(tag))
	}
	return tags
}
