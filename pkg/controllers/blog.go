package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlogController struct{}

func NewBlogController() *BlogController {
	return &BlogController{}
}

func (BlogController) ServeHome(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"message": "it is what it is!",
	})
}
