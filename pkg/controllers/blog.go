package controllers

import (
	"net/http"

	"github.com/Muhammed-Rajab/go-blog/pkg/db"
	"github.com/Muhammed-Rajab/go-blog/pkg/models"
	"github.com/gin-gonic/gin"
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

func (BlogController) BlogHandler(ctx *gin.Context) {

	obj := gin.H{}
	slug := ctx.Param("slug")
	blogs := models.NewBlogs(db.GetMDB().BlogsCollection())

	post, err := blogs.FindBlogBySlug(slug)

	if err != nil {
		obj["errors"] = err
	} else {
		obj["post"] = post
	}
	ctx.HTML(http.StatusOK, "blog.html", obj)
}
