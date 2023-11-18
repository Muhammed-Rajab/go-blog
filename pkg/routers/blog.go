package routers

import (
	"github.com/Muhammed-Rajab/go-blog/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func BlogRouter(root *gin.RouterGroup) *gin.RouterGroup {

	router := root.Group("/blog")
	controller := controllers.NewBlogController()
	{

		router.GET("", controller.HomeHandler)
		router.GET("/:slug", controller.BlogHandler)
	}

	return router
}
