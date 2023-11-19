package routers

import (
	"github.com/Muhammed-Rajab/go-blog/pkg/controllers"
	"github.com/Muhammed-Rajab/go-blog/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func BlogRouter(root *gin.RouterGroup) *gin.RouterGroup {

	router := root.Group("/blog")
	controller := controllers.NewBlogController()
	middlewares := middlewares.NewBlogMiddlewares()
	AuthMiddleware := middlewares.CheckForDashboardKey()
	{
		dashboard := router.Group("/dashboard")
		{
			dashboard.GET("", AuthMiddleware, controller.DashboardHandler)
			dashboard.DELETE("/:id", AuthMiddleware, controller.DashboardHandler)
		}

		router.GET("", controller.HomeHandler)
		router.GET("/:slug", controller.BlogHandler)

	}

	return router
}
