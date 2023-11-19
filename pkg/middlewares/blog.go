package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type BlogMiddlewares struct{}

func NewBlogMiddlewares() BlogMiddlewares {
	return BlogMiddlewares{}
}

func (BlogMiddlewares) CheckForDashboardKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		key, err := ctx.Cookie("auth-token")

		// Later, change this to rendering "Unauthorized" page
		if err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/blog")
		}

		if key == os.Getenv("BLOG_DASHBOARD_KEY") {
			ctx.Next()
		} else {
			ctx.Redirect(http.StatusTemporaryRedirect, "/blog")
		}
	}
}
