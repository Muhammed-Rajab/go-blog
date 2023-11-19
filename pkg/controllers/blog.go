package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
		// pass errors to the object and show the page, later
		// ctx.HTML(http.StatusOK, "add_blog.html", obj)
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

	oid, err := blogs.AddBlog(blog)
	if err != nil {
		obj["error"] = "Shit took a turn for the worst: " + err.Error()
		ctx.JSON(http.StatusBadRequest, obj)
		return
		// pass errors to the object and show the page, later
		// ctx.HTML(http.StatusOK, "add_blog.html", obj)
	}

	// Fetch the newly saved blog from db
	newBlog, err := blogs.FindBlogByID(oid.Hex())
	if err != nil {
		obj["error"] = "Shit took a turn for the worst: " + err.Error()
		ctx.JSON(http.StatusBadRequest, obj)
		return
		// pass errors to the object and show the page, later
		// ctx.HTML(http.StatusOK, "add_blog.html", obj)
	}

	// Redirect to the created blog if everything went well
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/blog/%s", newBlog.Slug))
}

func tagsFromString(stringTag string) []string {
	tags := strings.Split(stringTag, ",")
	for i, tag := range tags {
		tags[i] = strings.ToLower(strings.TrimSpace(tag))
	}
	return tags
}

func (BlogController) AuthDashboardHandler(ctx *gin.Context) {
	var obj gin.H
	ctx.HTML(http.StatusOK, "auth.html", obj)
}

func (BlogController) AuthDashboard(ctx *gin.Context) {
	form := struct {
		Token string `form:"token"`
	}{}

	if err := ctx.ShouldBind(&form); err != nil {
		// handle the token not found stuff
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "fuck off my property, mate",
		})
		return
	}

	log.Print(form.Token)

	// if the token is similar to the env variable,
	// then set it as cookie
	if form.Token == os.Getenv("BLOG_DASHBOARD_KEY") {
		ctx.SetCookie("auth-token", form.Token, 86400, "/", "localhost", false, true)
		ctx.Redirect(http.StatusSeeOther, "/blog/dashboard")
		return
	}
	log.Print("shit didn't work")

	// else render the auth page with error
	// later....
	ctx.Redirect(http.StatusSeeOther, "/blog/dashboard/auth")
}
