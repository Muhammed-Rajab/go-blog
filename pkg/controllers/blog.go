package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/Muhammed-Rajab/go-blog/pkg/db"
	"github.com/Muhammed-Rajab/go-blog/pkg/models"
	"github.com/Muhammed-Rajab/go-blog/pkg/utils"
	"github.com/gin-gonic/gin"
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
	// Add feature to sort by order date

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
		obj["error"] = err
	} else {
		obj["posts"] = posts
	}
	ctx.HTML(http.StatusOK, "index.html", obj)
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
		obj["error"] = err.Error()
		ctx.HTML(http.StatusBadRequest, "blog.html", obj)
		return
	} else {
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
	// Add feature to sort by date

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
		obj["error"] = err.Error()
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

func (BlogController) EditBlogHandler(ctx *gin.Context) {
	var obj gin.H

	id := ctx.Param("id")
	blogs := models.NewBlogs(db.GetMDB().BlogsCollection())

	blog, err := blogs.FindBlogByID(id)
	if err != nil {
		obj = gin.H{
			"error": "blog does not exists",
		}
		ctx.HTML(http.StatusBadRequest, "edit_blog.html", obj)
		return
	}

	publish := "off"
	if blog.Published {
		publish = "on"
	}

	form := models.BlogForm{
		ID:      blog.ID.Hex(),
		Title:   blog.Title,
		Desc:    blog.Desc,
		Content: blog.Content,
		Tags:    strings.Join(blog.Tags, ","),
		Publish: publish,
	}

	obj = gin.H{
		"post": form,
	}

	ctx.HTML(http.StatusOK, "edit_blog.html", obj)
}

func (BlogController) EditBlog(ctx *gin.Context) {
	var obj gin.H
	var form models.BlogForm

	if err := ctx.ShouldBind(&form); err != nil {
		obj = gin.H{
			"error": err.Error(),
		}
		ctx.HTML(http.StatusBadRequest, "edit_blog.html", obj)
		return
	}

	blogs := models.NewBlogs(db.GetMDB().BlogsCollection())

	// Create BlogModel from form
	publish := false
	if form.Publish == "on" {
		publish = true
	}

	oid, err := primitive.ObjectIDFromHex(form.ID)
	if err != nil {
		obj = gin.H{
			"error": err.Error(),
		}
		ctx.HTML(http.StatusBadRequest, "edit_blog.html", obj)
		return
	}

	// Get created time from object
	oldBlog, err := blogs.FindBlogByID(oid.Hex())
	if err != nil {
		obj = gin.H{
			"error": err.Error(),
		}
		ctx.HTML(http.StatusBadRequest, "edit_blog.html", obj)
		return
	}

	blog := models.BlogModel{
		ID:        oid,
		Title:     form.Title,
		Desc:      form.Desc,
		Content:   form.Content,
		Tags:      utils.TagsFromString(form.Tags),
		Published: publish,
		CreatedAt: oldBlog.CreatedAt,
	}

	err = blogs.UpdateBlogByID(oid.Hex(), blog)
	if err != nil {
		obj = gin.H{
			"error": err.Error(),
		}
		ctx.HTML(http.StatusBadRequest, "edit_blog.html", obj)
		return
	}

	newBlog, err := blogs.FindBlogByID(oid.Hex())
	if err != nil {
		obj = gin.H{
			"error": err.Error(),
		}
		ctx.HTML(http.StatusBadRequest, "edit_blog.html", obj)
		return
	}

	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/blog/%s", newBlog.Slug))
}

func (BlogController) AddBlogHandler(ctx *gin.Context) {
	var obj gin.H
	ctx.HTML(http.StatusOK, "add_blog.html", obj)
}

func (BlogController) AddBlog(ctx *gin.Context) {
	var obj gin.H
	var form models.BlogForm

	if err := ctx.ShouldBind(&form); err != nil {
		obj = gin.H{
			"error": err.Error(),
		}
		ctx.HTML(http.StatusBadRequest, "add_blog.html", obj)
		return
	}

	blogs := models.NewBlogs(db.GetMDB().BlogsCollection())

	publish := false
	if form.Publish == "on" {
		publish = true
	}

	blog := models.BlogModel{
		Title:     form.Title,
		Desc:      form.Desc,
		Content:   form.Content,
		Tags:      utils.TagsFromString(form.Tags),
		Published: publish,
	}

	oid, err := blogs.AddBlog(blog)
	if err != nil {
		obj = gin.H{
			"error": err.Error(),
		}
		ctx.HTML(http.StatusBadRequest, "add_blog.html", obj)
		return
	}

	newBlog, err := blogs.FindBlogByID(oid.Hex())
	if err != nil {
		obj = gin.H{
			"error": err.Error(),
		}
		ctx.HTML(http.StatusBadRequest, "add_blog.html", obj)
		return
	}

	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/blog/%s", newBlog.Slug))
}

func (BlogController) AuthDashboardHandler(ctx *gin.Context) {
	var obj gin.H
	ctx.HTML(http.StatusOK, "auth.html", obj)
}

func (BlogController) AuthDashboard(ctx *gin.Context) {
	form := struct {
		Token string `form:"token"`
	}{}
	var obj gin.H

	if err := ctx.ShouldBind(&form); err != nil {
		obj = gin.H{
			"error": "invalid token provided",
		}
		ctx.HTML(http.StatusUnauthorized, "auth.html", obj)
		return
	}

	// if the token is similar to the env variable,
	// then set it as cookie
	if form.Token == os.Getenv("BLOG_DASHBOARD_KEY") {
		ctx.SetCookie("auth-token", form.Token, 86400, "/", "localhost", false, true)
		ctx.Redirect(http.StatusSeeOther, "/blog/dashboard")
		return
	}

	// else render the auth page with error
	obj = gin.H{
		"error": "token does not match",
	}
	ctx.HTML(http.StatusBadRequest, "auth.html", obj)
}

func (BlogController) LogoutDashboard(ctx *gin.Context) {
	// expire auth-token cookie to logout
	ctx.SetCookie("auth-token", "", -1, "/", "localhost", false, true)
	ctx.Redirect(http.StatusSeeOther, "/blog/dashboard/auth")
}

func (BlogController) ImagesHandler(ctx *gin.Context) {
	obj := gin.H{}
	images := models.NewImages(db.GetMDB().ImagesCollection())

	pageNo := ctx.Query("page")
	search := ctx.Query("search")

	page, err := strconv.ParseUint(pageNo, 10, 64)
	if err != nil {
		page = 1
	}

	posts, err := images.FindImages(bson.M{
		"caption": bson.M{
			"$regex":   search,
			"$options": "i",
		},
	}, int(page), 10)

	if err != nil {
		obj["error"] = err
	} else {
		obj["images"] = posts
	}
	ctx.HTML(http.StatusOK, "images.html", obj)
}

func (BlogController) UploadImages(ctx *gin.Context) {

	var obj gin.H

	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
		obj = gin.H{
			"error": "failed to parse image. " + err.Error(),
		}
		ctx.HTML(http.StatusBadRequest, "images.html", obj)
		return
	}

	caption := ctx.Request.PostFormValue("caption")
	if caption == "" {
		obj = gin.H{
			"error": "empty caption not allowed",
		}
		ctx.HTML(http.StatusBadRequest, "images.html", obj)
		return
	}

	file, _, err := ctx.Request.FormFile("image")
	if err != nil {
		obj = gin.H{
			"error": "error retrieving the file. " + err.Error(),
		}
		ctx.HTML(http.StatusBadRequest, "images.html", obj)
		return
	}
	defer file.Close()

	// Move this to .env later
	uploadDir := "./public/uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	// Create slug kinda stuff for the filename from the caption
	images := models.NewImages(db.GetMDB().ImagesCollection())
	slug := images.CreateSlug(caption)
	filename := filepath.Join(uploadDir, slug)

	out, err := os.Create(filename)
	if err != nil {
		obj = gin.H{
			"error": "error saving the file. " + err.Error(),
		}
		ctx.HTML(http.StatusInternalServerError, "images.html", obj)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		obj = gin.H{
			"error": "error copying the file. " + err.Error(),
		}
		ctx.HTML(http.StatusInternalServerError, "images.html", obj)
		return
	}

	// If things are successful, then save image to database
	image := models.ImageModel{
		Caption:  caption,
		Location: filename,
		Slug:     slug,
	}

	if _, err := images.AddImage(image); err != nil {
		obj = gin.H{
			"error": "error saving the image to database: " + err.Error(),
		}
		ctx.HTML(http.StatusInternalServerError, "images.html", obj)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/blog/dashboard/images")
}

func (BlogController) DeleteImage(ctx *gin.Context) {
	obj := gin.H{}
	images := models.NewImages(db.GetMDB().ImagesCollection())

	id := ctx.Param("id")

	if err := images.DeleteImageByID(id); err != nil {
		obj["error"] = err.Error()
		ctx.HTML(http.StatusBadRequest, "images.html", obj)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/blog/dashboard/images")
}
